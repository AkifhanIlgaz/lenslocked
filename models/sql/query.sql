SELECT users.id,
    users.email,
    users.password_hash
FROM sessions
    JOIN users ON users.id = sessions.user_id
WHERE sessions.token_hash = '0M3jdFRpoBY6fQoRvDSQzZbgNwYnI4dLEGUtLsxnSBA=';