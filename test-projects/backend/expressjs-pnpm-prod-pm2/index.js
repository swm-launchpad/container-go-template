import express from 'express';

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(express.json());

// In-memory data store
let items = [
  { id: 1, name: 'Sample Item 1', description: 'First sample item', createdAt: new Date().toISOString() },
  { id: 2, name: 'Sample Item 2', description: 'Second sample item', createdAt: new Date().toISOString() }
];
let nextId = 3;
let requestCount = 0;
const startTime = Date.now();

// Request counter middleware
app.use((req, res, next) => {
  requestCount++;
  next();
});

// GET / - Root endpoint with API info
app.get('/', (req, res) => {
  res.json({
    message: 'Express.js Test API',
    version: '1.0.0',
    requestCount,
    timestamp: new Date().toISOString(),
    endpoints: [
      'GET /',
      'GET /health',
      'GET /items',
      'POST /items',
      'GET /items/:id',
      'DELETE /items/:id'
    ]
  });
});

// GET /health - Health check
app.get('/health', (req, res) => {
  const uptime = Math.floor((Date.now() - startTime) / 1000);
  res.json({
    status: 'healthy',
    uptime: `${uptime} seconds`,
    timestamp: new Date().toISOString()
  });
});

// GET /items - List all items
app.get('/items', (req, res) => {
  res.json({
    success: true,
    count: items.length,
    data: items
  });
});

// POST /items - Create new item
app.post('/items', (req, res) => {
  const { name, description } = req.body;

  if (!name) {
    return res.status(400).json({
      success: false,
      error: 'Name is required'
    });
  }

  const newItem = {
    id: nextId++,
    name,
    description: description || '',
    createdAt: new Date().toISOString()
  };

  items.push(newItem);

  res.status(201).json({
    success: true,
    data: newItem
  });
});

// GET /items/:id - Get specific item
app.get('/items/:id', (req, res) => {
  const id = parseInt(req.params.id);
  const item = items.find(i => i.id === id);

  if (!item) {
    return res.status(404).json({
      success: false,
      error: 'Item not found'
    });
  }

  res.json({
    success: true,
    data: item
  });
});

// DELETE /items/:id - Delete item
app.delete('/items/:id', (req, res) => {
  const id = parseInt(req.params.id);
  const index = items.findIndex(i => i.id === id);

  if (index === -1) {
    return res.status(404).json({
      success: false,
      error: 'Item not found'
    });
  }

  const deletedItem = items.splice(index, 1)[0];

  res.json({
    success: true,
    message: 'Item deleted successfully',
    data: deletedItem
  });
});

// 404 handler
app.use((req, res) => {
  res.status(404).json({
    success: false,
    error: 'Endpoint not found'
  });
});

// Error handler
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    success: false,
    error: 'Internal server error'
  });
});

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
