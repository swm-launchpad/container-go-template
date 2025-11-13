// MongoDB test database initialization
// This script creates comprehensive test data for validating MongoDB functionality

db = db.getSiblingDB('testdb');

// Drop existing collections to ensure clean state
db.users.drop();
db.products.drop();
db.orders.drop();

// Insert test users with diverse MongoDB document structures
db.users.insertMany([
    {
        username: 'john_doe',
        email: 'john@example.com',
        age: 28,
        balance: 1500.50,
        isActive: true,
        preferences: {
            theme: 'dark',
            notifications: true,
            language: 'en'
        },
        tags: ['premium', 'verified'],
        createdAt: new Date('2024-01-15'),
        updatedAt: new Date()
    },
    {
        username: 'jane_smith',
        email: 'jane@example.com',
        age: 34,
        balance: 2750.00,
        isActive: true,
        preferences: {
            theme: 'light',
            notifications: false,
            language: 'en'
        },
        tags: ['premium', 'verified', 'early-adopter'],
        createdAt: new Date('2024-01-10'),
        updatedAt: new Date()
    },
    {
        username: 'bob_wilson',
        email: 'bob@example.com',
        age: 45,
        balance: 500.25,
        isActive: false,
        preferences: {
            theme: 'dark',
            notifications: true,
            language: 'es'
        },
        tags: ['verified'],
        createdAt: new Date('2024-02-20'),
        updatedAt: new Date()
    },
    {
        username: 'alice_brown',
        email: 'alice@example.com',
        age: 22,
        balance: 3200.75,
        isActive: true,
        preferences: {
            theme: 'light',
            notifications: true,
            language: 'fr'
        },
        tags: ['premium', 'verified', 'influencer'],
        createdAt: new Date('2024-01-05'),
        updatedAt: new Date()
    },
    {
        username: 'charlie_davis',
        email: 'charlie@example.com',
        age: 31,
        balance: 0.00,
        isActive: true,
        preferences: {
            theme: 'auto',
            notifications: false,
            language: 'en'
        },
        tags: [],
        createdAt: new Date('2024-03-01'),
        updatedAt: new Date()
    }
]);

// Insert test products with nested documents and arrays
db.products.insertMany([
    {
        name: 'Laptop Pro 15',
        description: 'High-performance laptop with 16GB RAM',
        price: 1299.99,
        stock: 25,
        category: 'Electronics',
        rating: 4.5,
        tags: ['laptop', 'computer', 'portable'],
        specifications: {
            brand: 'TechCorp',
            warranty: '2 years',
            color: 'silver',
            ram: '16GB',
            storage: '512GB SSD'
        },
        reviews: [
            { user: 'john_doe', rating: 5, comment: 'Excellent laptop!', date: new Date('2024-02-01') },
            { user: 'jane_smith', rating: 4, comment: 'Great performance', date: new Date('2024-02-05') }
        ],
        createdAt: new Date('2024-01-01')
    },
    {
        name: 'Wireless Mouse',
        description: 'Ergonomic wireless mouse with USB receiver',
        price: 29.99,
        stock: 150,
        category: 'Electronics',
        rating: 4.2,
        tags: ['mouse', 'wireless', 'accessory'],
        specifications: {
            brand: 'PeripheralPlus',
            battery: 'AA',
            color: 'black',
            dpi: '1600'
        },
        reviews: [],
        createdAt: new Date('2024-01-15')
    },
    {
        name: 'Coffee Maker',
        description: 'Programmable 12-cup coffee maker',
        price: 79.99,
        stock: 50,
        category: 'Home & Kitchen',
        rating: 4.0,
        tags: ['coffee', 'kitchen', 'appliance'],
        specifications: {
            brand: 'BrewMaster',
            capacity: '12 cups',
            color: 'stainless',
            programmable: true
        },
        reviews: [
            { user: 'bob_wilson', rating: 4, comment: 'Makes good coffee', date: new Date('2024-02-10') }
        ],
        createdAt: new Date('2024-01-10')
    },
    {
        name: 'Running Shoes',
        description: 'Comfortable running shoes for all terrains',
        price: 89.99,
        stock: 75,
        category: 'Sports',
        rating: 4.7,
        tags: ['shoes', 'running', 'sports'],
        specifications: {
            brand: 'RunFast',
            sizes: ['7', '8', '9', '10', '11'],
            color: 'blue',
            waterproof: true
        },
        reviews: [
            { user: 'alice_brown', rating: 5, comment: 'Very comfortable!', date: new Date('2024-02-15') }
        ],
        createdAt: new Date('2024-01-20')
    },
    {
        name: 'Desk Lamp',
        description: 'LED desk lamp with adjustable brightness',
        price: 45.99,
        stock: 100,
        category: 'Home & Kitchen',
        rating: 4.3,
        tags: ['lamp', 'led', 'desk'],
        specifications: {
            brand: 'BrightLight',
            wattage: '10W',
            color: 'white',
            adjustable: true
        },
        reviews: [],
        createdAt: new Date('2024-01-25')
    },
    {
        name: 'Backpack',
        description: 'Water-resistant laptop backpack',
        price: 59.99,
        stock: 60,
        category: 'Accessories',
        rating: 4.6,
        tags: ['backpack', 'travel', 'storage'],
        specifications: {
            brand: 'PackIt',
            capacity: '30L',
            color: 'navy',
            waterResistant: true
        },
        reviews: [
            { user: 'charlie_davis', rating: 5, comment: 'Perfect for travel', date: new Date('2024-03-05') }
        ],
        createdAt: new Date('2024-02-01')
    }
]);

// Insert test orders with references and embedded data
db.orders.insertMany([
    {
        userId: 'john_doe',
        items: [
            { productName: 'Laptop Pro 15', quantity: 1, price: 1299.99 }
        ],
        totalPrice: 1299.99,
        status: 'completed',
        shippingAddress: {
            street: '123 Main St',
            city: 'New York',
            state: 'NY',
            zipCode: '10001',
            country: 'USA'
        },
        orderDate: new Date('2024-02-01'),
        deliveryDate: new Date('2024-02-05')
    },
    {
        userId: 'john_doe',
        items: [
            { productName: 'Wireless Mouse', quantity: 2, price: 29.99 }
        ],
        totalPrice: 59.98,
        status: 'completed',
        shippingAddress: {
            street: '123 Main St',
            city: 'New York',
            state: 'NY',
            zipCode: '10001',
            country: 'USA'
        },
        orderDate: new Date('2024-02-10'),
        deliveryDate: new Date('2024-02-13')
    },
    {
        userId: 'jane_smith',
        items: [
            { productName: 'Coffee Maker', quantity: 1, price: 79.99 }
        ],
        totalPrice: 79.99,
        status: 'processing',
        shippingAddress: {
            street: '456 Oak Ave',
            city: 'Los Angeles',
            state: 'CA',
            zipCode: '90001',
            country: 'USA'
        },
        orderDate: new Date('2024-03-01'),
        deliveryDate: null
    },
    {
        userId: 'alice_brown',
        items: [
            { productName: 'Running Shoes', quantity: 1, price: 89.99 },
            { productName: 'Backpack', quantity: 1, price: 59.99 }
        ],
        totalPrice: 149.98,
        status: 'completed',
        shippingAddress: {
            street: '789 Pine Rd',
            city: 'Chicago',
            state: 'IL',
            zipCode: '60601',
            country: 'USA'
        },
        orderDate: new Date('2024-02-20'),
        deliveryDate: new Date('2024-02-25')
    },
    {
        userId: 'charlie_davis',
        items: [
            { productName: 'Desk Lamp', quantity: 2, price: 45.99 }
        ],
        totalPrice: 91.98,
        status: 'pending',
        shippingAddress: {
            street: '321 Elm St',
            city: 'Houston',
            state: 'TX',
            zipCode: '77001',
            country: 'USA'
        },
        orderDate: new Date('2024-03-10'),
        deliveryDate: null
    }
]);

// Create indexes for better query performance
db.users.createIndex({ username: 1 }, { unique: true });
db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ isActive: 1 });
db.products.createIndex({ category: 1 });
db.products.createIndex({ tags: 1 });
db.orders.createIndex({ userId: 1 });
db.orders.createIndex({ status: 1 });
db.orders.createIndex({ orderDate: -1 });

// Display summary statistics
print('MongoDB initialized successfully!');
print('-----------------------------------');
print('Total users:', db.users.countDocuments());
print('Total products:', db.products.countDocuments());
print('Total orders:', db.orders.countDocuments());
print('Active users:', db.users.countDocuments({ isActive: true }));
print('Products in Electronics category:', db.products.countDocuments({ category: 'Electronics' }));
print('Completed orders:', db.orders.countDocuments({ status: 'completed' }));
print('-----------------------------------');
