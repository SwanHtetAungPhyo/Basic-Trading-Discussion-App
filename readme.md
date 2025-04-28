## Basic Trading Discussion App Design 

### 1. Overview
This trading app will provide real-time market data for Solana (SOL) trading pairs, specifically through Binance's WebSocket API, and support core trading operations: viewing live prices, placing market and limit orders, and tracking order status and trade history.

### 2. Key Features
1. **Real‑Time Price Chart**
    - 1-minute candlestick chart for SOL/USDT
    - Live price ticker updates

2. **Order Entry**
    - Market orders (buy/sell at current market price)
    - Limit orders (buy/sell at specified price)
    - Simple order form with inputs for quantity and price

3. **Order Management**
    - View open orders
    - Cancel orders
    - View trade history (past executions)

4. **Account Dashboard**
    - Display account balances for SOL and USDT
    - Show P&L for open positions

5. **Notifications**
    - Order execution and cancellation confirmations
    - Price alerts (optional)

6. **Technical Indicators**
    - Support for MACD:
        - **Fast Line** (default 12-period EMA)
        - **Slow Line** (default 26-period EMA)
        - **Signal Line** (default 9-period EMA of MACD)
    - Other indicators (SMA, RSI, EMA) already supported

---

## Documentation

**Real-time SOL price and technical indicators and user discussion backend services**

### Overview

Users can see real-time updates of the Solana price on their dashboard with:

- EMA
- SMA
- RSI
- Candlestick chart

Users can discuss trades via chat rooms.

### Technology Architecture

```
User → Webpage → Login (Lambda Service) → Dashboard (Lambda Function) → Real-time Data (EC2 Instance)
                               ↑                     ↓
                             DynamoDB (CRUD & Lambda) ← SQS Queue ← Data feed from Binance & aggregation
```

**Technology Stack:**

- Go (Fiber)
- AWS Simple Queue Service (SQS)
- AWS DynamoDB
- AWS EC2
- AWS Lambda
- Docker

