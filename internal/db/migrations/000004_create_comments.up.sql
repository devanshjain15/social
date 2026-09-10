CREATE TABLE COMMENTS (
    id bigserial PRIMARY KEY, 
    content text NOT NULL, 
    post_id bigserial NOT NULL, 
    user_id bigserial NOT NULL, 
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
)
  