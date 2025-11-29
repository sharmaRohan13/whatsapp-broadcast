#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';
import csv from 'csv-parser';
import pkg from 'whatsapp-web.js';
const { Client, LocalAuth } = pkg;
import qrcode from 'qrcode-terminal';
import { Command } from 'commander';
import chalk from 'chalk';

// Get the directory where this script is located
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// --- Utility: Read CSV ---
const readContacts = (filePath) => {
  return new Promise((resolve, reject) => {
    const contacts = [];
    fs.createReadStream(filePath)
      .pipe(csv())
      .on('data', (row) => {
        // Handle CSV with spaces in headers by checking all possible variations
        const name = (row.name || row[' name'] || row['name '] || row[' name '])?.trim();
        const number = (row.number || row[' number'] || row['number '] || row[' number '] || row[' number,'])?.trim().replace(/,$/, '');
        
        if (name && number) {
          contacts.push({ name, number });
        }
      })
      .on('end', () => resolve(contacts))
      .on('error', reject);
  });
};

// --- Utility: Random delay to avoid detection ---
const sleep = (min, max) =>
  new Promise((r) => setTimeout(r, Math.floor(Math.random() * (max - min + 1)) + min));

// --- CLI Setup ---
const program = new Command();

program
  .name('whatsapp-broadcast')
  .description('Send WhatsApp broadcast messages to multiple contacts')
  .version('1.0.0')
  .option('-n, --numbers <path>', 'Path to CSV file with contacts (name,number)', '../sample/numbers.csv')
  .option('-m, --message <path>', 'Path to text file with message template', '../sample/message.txt')
  .option('-f, --full', 'Send to all contacts (default: test mode with first contact only)')
  .option('-d, --delay <range>', 'Delay range in seconds (e.g., "15-35")', '15-35')
  .parse(process.argv);

const options = program.opts();

// Parse delay range
const [minDelay, maxDelay] = options.delay.split('-').map(n => parseInt(n) * 1000);
if (isNaN(minDelay) || isNaN(maxDelay)) {
  console.error(chalk.red('❌ Invalid delay format. Use format: 15-35'));
  process.exit(1);
}

const numbersPath = path.resolve(options.numbers);
const messagePath = path.resolve(options.message);
const isTestMode = !options.full;

// Validate file existences
if (!fs.existsSync(numbersPath)) {
  console.error(chalk.red(`❌ Numbers file not found: ${numbersPath}`));
  process.exit(1);
}

if (!fs.existsSync(messagePath)) {
  console.error(chalk.red(`❌ Message file not found: ${messagePath}`));
  process.exit(1);
}

console.log(`dataPath: ${path.join(__dirname, '.wwebjs_auth')}`);
// --- WhatsApp client setup ---
const client = new Client({
  authStrategy: new LocalAuth({
    dataPath: path.join(__dirname, '.wwebjs_auth'),
    clientId: 'whatsapp-node-client'
  }),
  puppeteer: {
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  }
});

// --- Generate QR for first-time login ---
client.on('qr', (qr) => {
  console.log(chalk.cyan('\n📱 Scan this QR code in WhatsApp (Menu → Linked Devices):'));
  qrcode.generate(qr, { small: true });
});

client.on('authenticated', () => {
  console.log(chalk.green('✅ Authenticated successfully'));
});

client.on('auth_failure', () => {
  console.error(chalk.red('❌ Authentication failed'));
  process.exit(1);
});

client.on('ready', async () => {
  console.log(chalk.green('✅ WhatsApp client ready.\n'));

  const messageTemplate = fs.readFileSync(messagePath, 'utf8').trim();
  const contacts = await readContacts(numbersPath);

  if (contacts.length === 0) {
    console.error(chalk.red('❌ No valid contacts found in CSV file'));
    process.exit(1);
  }

  // Use test mode by default (first contact only), unless --full flag is passed
  const testContacts = isTestMode ? contacts.slice(0, 1) : contacts;
  const modeText = isTestMode ? chalk.yellow('TEST MODE - first contact only') : chalk.green('FULL MODE - all contacts');

  console.log(chalk.blue(`📋 Found ${contacts.length} contacts, running in ${modeText}`));
  console.log(chalk.blue(`⏱️  Delay between messages: ${minDelay/1000}-${maxDelay/1000} seconds`));
  console.log(chalk.gray(`\nMessage template:\n"${messageTemplate.substring(0, 100)}${messageTemplate.length > 100 ? '...' : ''}"\n`));
  console.log(chalk.cyan(`🚀 Starting ${isTestMode ? 'test' : 'full'} broadcast...\n`));

  let successCount = 0;
  let failCount = 0;

  for (const [i, contact] of testContacts.entries()) {
    const { name, number } = contact;
    const cleanNum = number.replace(/[^\d]/g, '');
    const chatId = `${cleanNum}@c.us`;
    
    // Replace ${name} placeholder with actual name
    const personalizedMessage = messageTemplate.replace(/\${name}/g, name);

    console.log(chalk.gray(`   Attempting to send to: ${chatId}`));

    try {
      // First check if the number exists on WhatsApp
      const numberId = await client.getNumberId(cleanNum);
      
      if (!numberId) {
        failCount++;
        console.error(chalk.yellow(`⚠️  [${i + 1}/${testContacts.length}] ${name} (${number}) is not registered on WhatsApp`));
        continue;
      }

      const sentMessage = await client.sendMessage(numberId._serialized, personalizedMessage);
      if (sentMessage && sentMessage.id) {
        successCount++;
        console.log(chalk.green(`✅ [${i + 1}/${testContacts.length}] Sent to ${name} (${number})`));
        console.log(chalk.gray(`   Message ID: ${sentMessage.id.id}`));
      } else {
        failCount++;
        console.error(chalk.yellow(`⚠️  [${i + 1}/${testContacts.length}] Sent but no confirmation for ${name} (${number})`));
      }
    } catch (err) {
      failCount++;
      console.error(chalk.red(`❌ [${i + 1}/${testContacts.length}] Failed for ${name} (${number})`));
      console.error(chalk.red(`   Error: ${err.message}`));
      if (err.stack) {
        console.error(chalk.gray(`   ${err.stack.split('\n')[1]}`));
      }
    }

    // Don't delay after the last message
    if (i < testContacts.length - 1) {
      const delayTime = Math.floor(Math.random() * (maxDelay - minDelay + 1)) + minDelay;
      console.log(chalk.gray(`   ⏳ Waiting ${Math.round(delayTime/1000)}s before next message...\n`));
      await sleep(minDelay, maxDelay);
    }
  }

  console.log(chalk.green(`\n🎉 ${isTestMode ? 'Test' : 'Full'} broadcast completed!`));
  console.log(chalk.blue(`📊 Summary: ${successCount} successful, ${failCount} failed`));
  process.exit(0);
});

client.initialize();
