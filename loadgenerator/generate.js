const { EventHubProducerClient } = require("@azure/event-hubs");
const { v4: uuidv4 } = require('uuid');

async function generate(messageCount) {
  messageCount = parseInt(messageCount);

  const producerClient = new EventHubProducerClient(process.env.EVENTHUB_CONNECTION);

  let eventDataBatch = await producerClient.createBatch();

  while (messageCount > 0) {
    let wasAdded = eventDataBatch.tryAdd({ body: uuidv4() });
    if (wasAdded) {
      messageCount -= 1;
    }
    if (!wasAdded || messageCount === 0) {
      console.log(`Sending ${eventDataBatch.count} messages, ${messageCount} remaining...`);
      await producerClient.sendBatch(eventDataBatch);
      eventDataBatch = await producerClient.createBatch();
    }
  }

  await producerClient.close();
}

const messageCount = process.argv[2] || 1;
generate(messageCount);