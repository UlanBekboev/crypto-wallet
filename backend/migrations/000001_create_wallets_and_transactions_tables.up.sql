CREATE TABLE wallets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  balance BIGINT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sender_wallet_id UUID REFERENCES wallets(id) ON DELETE SET NULL,
  receiver_wallet_id UUID REFERENCES wallets(id) ON DELETE SET NULL,
  amount BIGINT NOT NULL CHECK (amount > 0),
  created_at TIMESTAMPTZ DEFAULT NOW()
);