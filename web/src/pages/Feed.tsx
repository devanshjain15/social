import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api, type Post } from '../lib/api';
import './Feed.css';

export function Feed() {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [hasMore, setHasMore] = useState(true);
  const limit = 10;
  const [currentOffset, setCurrentOffset] = useState(0);

  useEffect(() => {
    loadPosts();
  }, []);

  const loadPosts = async () => {
    try {
      const posts = await api.getUserFeed(limit, currentOffset);
      if (posts.length < limit) {
        setHasMore(false);
      }
      setPosts(posts);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load feed');
    } finally {
      setLoading(false);
    }
  };

  const loadMore = async () => {
    const newOffset = currentOffset + limit;
    try {
      const morePosts = await api.getUserFeed(limit, newOffset);
      if (morePosts.length < limit) {
        setHasMore(false);
      }
      setPosts([...posts, ...morePosts]);
      setCurrentOffset(newOffset);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load more posts');
    }
  };

  if (loading) {
    return (
      <div className="feed-container">
        <div className="feed-loading">
          <div className="spinner"></div>
          <p>Loading your feed...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="feed-container">
        <div className="feed-error">
          <p>{error}</p>
          <button onClick={loadPosts} className="retry-button">
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="feed-container">
      <div className="feed-header">
        <h1>Your Feed</h1>
        <Link to="/create-post" className="create-post-button">
          + New Post
        </Link>
      </div>

      {posts.length === 0 ? (
        <div className="feed-empty">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
            <path d="M19 3H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z" stroke="currentColor" strokeWidth="2"/>
          </svg>
          <h2>Your feed is empty</h2>
          <p>Follow users to see their posts here</p>
          <Link to="/explore" className="explore-link">
            Discover Users
          </Link>
        </div>
      ) : (
        <>
          <div className="posts-list">
            {posts.map((post) => (
              <article key={post.id} className="post-card">
                <div className="post-header">
                  <div className="post-author">
                    <div className="post-avatar">
                      {post.user?.username?.[0]?.toUpperCase() || 'U'}
                    </div>
                    <div>
                      <div className="post-username">{post.user?.username || 'Unknown'}</div>
                      <div className="post-date">
                        {new Date(post.created_at).toLocaleDateString()}
                      </div>
                    </div>
                  </div>
                </div>

                <Link to={`/posts/${post.id}`} className="post-link">
                  <h2 className="post-title">{post.title}</h2>
                  <p className="post-content">{post.content}</p>
                </Link>

                {post.tags && post.tags.length > 0 && (
                  <div className="post-tags">
                    {post.tags.map((tag, index) => (
                      <span key={index} className="post-tag">
                        #{tag}
                      </span>
                    ))}
                  </div>
                )}

                <div className="post-footer">
                  <span className="post-comments">
                    💬 {post.comments_count || 0} comments
                  </span>
                </div>
              </article>
            ))}
          </div>

          {hasMore && (
            <button onClick={loadMore} className="load-more-button">
              Load More
            </button>
          )}
        </>
      )}
    </div>
  );
}
