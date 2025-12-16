-- Create user_pptk pivot table for many-to-many relationship
CREATE TABLE user_pptk (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pptk_id UUID NOT NULL REFERENCES pptk(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, pptk_id)
);

CREATE INDEX idx_user_pptk_user_id ON user_pptk(user_id);
CREATE INDEX idx_user_pptk_pptk_id ON user_pptk(pptk_id);
