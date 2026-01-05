-- Migration script to update existing businesses with owner_id
-- Run this AFTER backend has successfully started and migrations have completed

-- Step 1: Find first user to assign businesses to
-- If multiple users exist, you may need to customize this
DO $$
DECLARE
    first_user_id UUID;
    business_count INTEGER;
BEGIN
    SELECT id INTO first_user_id FROM users LIMIT 1;

    IF first_user_id IS NULL THEN
        RAISE NOTICE 'No users found. No migration needed.';
    ELSE
        -- Update all businesses with NULL owner_id to first user
        UPDATE businesses
        SET owner_id = first_user_id
        WHERE owner_id IS NULL;

        GET DIAGNOSTICS business_count = ROW_COUNT;
        RAISE NOTICE 'Updated % businesses with owner_id = %', business_count, first_user_id;
    END IF;
END $$;

-- Step 2: Verify migration
SELECT
    COUNT(*) as total_businesses,
    COUNT(owner_id) as businesses_with_owner,
    COUNT(*) - COUNT(owner_id) as businesses_without_owner
FROM businesses;
