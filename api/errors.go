package api

// General errors
const ErrInternal int32 = 500
const ErrInitialConfig int32 = 980

// Article errors
const ErrUnmarshalArticle int32 = 1001
const ErrNoCollection int32 = 1002
const ErrArticleNotFound int32 = 1003

// User errors
const ErrCreateOwner int32 = 2001
