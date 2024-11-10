-- Users table
CREATE TABLE IF NOT EXISTS Users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(100) UNIQUE,
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Loans table
CREATE TABLE IF NOT EXISTS Loans (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES Users(id) ON DELETE CASCADE,
    principal INT NOT NULL,
    interest_rate DECIMAL(5, 2) NOT NULL,
    total_weeks INT NOT NULL,
    weekly_amount INT NOT NULL,
    outstanding INT NOT NULL,
    missed_amount INT DEFAULT 0,
    delinquent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Payments table
CREATE TABLE IF NOT EXISTS Payments (
    id SERIAL PRIMARY KEY,
    loan_id INT REFERENCES Loans(id) ON DELETE CASCADE,
    amount_paid INT NOT NULL,
    payment_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Loan Schedule table
CREATE TABLE IF NOT EXISTS Loan_Schedules (
    id SERIAL PRIMARY KEY,
    loan_id INT REFERENCES Loans(id) ON DELETE CASCADE,
    week_number INT NOT NULL,
    amount_due INT NOT NULL,
    due_date DATE NOT NULL,
    is_paid BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Users table index
CREATE INDEX IF NOT EXISTS idx_users_username ON Users(username);

-- Loans table indexes
CREATE INDEX IF NOT EXISTS idx_loans_user_id ON Loans(user_id);

-- Payments table indexes
CREATE INDEX IF NOT EXISTS idx_payments_loan_id ON Payments(loan_id);

-- Loan_Schedules table indexes
CREATE INDEX IF NOT EXISTS idx_loan_schedules_loan_id ON Loan_Schedules(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_schedules_due_date ON Loan_Schedules(due_date);
CREATE INDEX IF NOT EXISTS idx_loan_schedules_is_paid ON Loan_Schedules(is_paid);

-- Combined index for overdue schedules query
CREATE INDEX IF NOT EXISTS idx_loan_schedules_loan_id_due_date_is_paid
ON Loan_Schedules(loan_id, due_date, is_paid);

-- Truncate all tables to start fresh
TRUNCATE TABLE Payments, Loan_Schedules, Loans, Users RESTART IDENTITY CASCADE;

-- Insert example user
INSERT INTO Users (name, username, email) VALUES ('Abia Lemuel', 'abialemuel', 'abialemuel@example.com');

-- Insert example loan for the user
-- Principal: Rp 5,000,000
-- Interest: 10% per annum flat rate (Rp 500,000 over the year, or Rp 10,000 per week for 50 weeks)
-- Total per week: Rp 110,000
INSERT INTO Loans (user_id, principal, interest_rate, total_weeks, weekly_amount, outstanding, missed_amount, delinquent) 
VALUES (1, 5000000, 10.0, 50, 110000, 5500000, 0, FALSE);

-- Seed Loan Schedules for the loan
-- Assuming we start on today’s date for the first due date, with each due weekly
DO
$$
DECLARE
    week_num INT := 1;
    due_date DATE := CURRENT_DATE;
    loan_id INT;
BEGIN
    -- Fetch the loan_id of the last inserted loan (assuming it's the only one)
    SELECT id INTO loan_id FROM Loans WHERE user_id = 1 ORDER BY id DESC LIMIT 1;
    
    WHILE week_num <= 50 LOOP
        INSERT INTO Loan_Schedules (loan_id, week_number, amount_due, due_date, is_paid) 
        VALUES (loan_id, week_num, 110000, due_date, FALSE);
        
        -- Increment to the next week
        week_num := week_num + 1;
        due_date := due_date + INTERVAL '1 week';
    END LOOP;
END
$$;

