// ============================================
// E-Wallet MongoDB Initialization Script
// ============================================

// Create collections
db.createCollection("notifications");
db.createCollection("ledger");
db.createCollection("analytics");

// Create indexes
db.notifications.createIndex({ "user_id": 1 });
db.notifications.createIndex({ "status": 1 });
db.notifications.createIndex({ "created_at": -1 });

db.ledger.createIndex({ "transaction_id": 1 });
db.ledger.createIndex({ "wallet_id": 1 });
db.ledger.createIndex({ "created_at": -1 });

db.analytics.createIndex({ "user_id": 1 });
db.analytics.createIndex({ "event": 1 });
db.analytics.createIndex({ "timestamp": -1 });

// Insert sample notification
db.notifications.insertOne({
    user_id: "550e8400-e29b-41d4-a716-446655440000",
    type: "email",
    title: "Welcome to E-Wallet!",
    content: "Thank you for registering. Your wallet has been created with IDR 1,000,000 initial balance.",
    status: "sent",
    created_at: new Date()
});

// Insert sample ledger entry
db.ledger.insertOne({
    transaction_id: "550e8400-e29b-41d4-a716-446655440010",
    wallet_id: "550e8400-e29b-41d4-a716-446655440001",
    type: "initial_balance",
    amount: 1000000,
    balance_before: 0,
    balance_after: 1000000,
    metadata: {
        description: "Initial wallet balance"
    },
    created_at: new Date()
});

print("✅ MongoDB initialized successfully!");
