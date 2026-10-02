import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { api, type Post } from '../lib/api';
import './PostDetail.css';

export function PostDetail() {
  const { postId } = useParams<{ postId: string }>();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [commentText, setCommentText] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    if (!postId) return;

    const loadPost = async () => {
      try {
        const data = await api.getPost(parseInt(postId));
        setPost(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load post');
      } finally {
        setLoading(false);
      }
    };

    loadPost();
  }, [postId]);

  const handleDelete = async () => {
    if (!post || !confirm('Are you sure you want to delete this post?')) {
      return;
    }

    try {
      await api.deletePost(post.id);
      navigate('/feed');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to delete post');
    }
  };

  const handleCommentSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!commentText.trim() || !post) return;

    setIsSubmitting(true);
    try {
      const newComment = await api.createComment(post.id, commentText);
      setPost({
        ...post,
        comments: [...(post.comments || []), newComment],
      });
      setCommentText('');
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to post comment');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="post-detail-container">
        <div className="post-detail-loading">
          <div className="spinner"></div>
          <p>Loading post...</p>
        </div>
      </div>
    );
  }

  if (error || !post) {
    return (
      <div className="post-detail-container">
        <div className="post-detail-error">
          <p>{error || 'Post not found'}</p>
          <Link to="/feed" className="back-button">
            Back to Feed
          </Link>
        </div>
      </div>
    );
  }

  const formatDate = (date: string) => {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className="post-detail-container">
      <div className="post-detail-nav">
        <Link to="/feed" className="back-link">
          ← Back to Feed
        </Link>
      </div>

      <article className="post-detail-card">
        <div className="post-detail-header">
          <div className="post-author">
            <div className="post-avatar">
              {post.user?.username?.[0]?.toUpperCase() || 'U'}
            </div>
            <div>
              <div className="post-username">{post.user?.username || 'Unknown'}</div>
              <div className="post-date">{formatDate(post.created_at)}</div>
            </div>
          </div>
          <button onClick={handleDelete} className="delete-button" aria-label="Delete post">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
          </button>
        </div>

        <h1 className="post-title">{post.title}</h1>

        <div className="post-content">{post.content}</div>

        {post.tags && post.tags.length > 0 && (
          <div className="post-tags">
            {post.tags.map((tag, index) => (
              <span key={index} className="post-tag">
                #{tag}
              </span>
            ))}
          </div>
        )}

        <div className="post-stats">
          <div className="post-stat">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
              <path d="M21 11.5a8.38 8.38 0 01-.9 3.8 8.5 8.5 0 01-7.6 4.7 8.38 8.38 0 01-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 01-.9-3.8 8.5 8.5 0 014.7-7.6 8.38 8.38 0 013.8-.9h.5a8.48 8.48 0 018 8v.5z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
            <span>{post.comments?.length || 0} comments</span>
          </div>
        </div>

        <div className="add-comment-section">
          <h2>Add a Comment</h2>
          <form onSubmit={handleCommentSubmit}>
            <textarea
              className="comment-input"
              placeholder="Write your comment..."
              value={commentText}
              onChange={(e) => setCommentText(e.target.value)}
              rows={3}
              disabled={isSubmitting}
            />
            <button
              type="submit"
              className="submit-comment-button"
              disabled={!commentText.trim() || isSubmitting}
            >
              {isSubmitting ? 'Posting...' : 'Post Comment'}
            </button>
          </form>
        </div>

        {post.comments && post.comments.length > 0 && (
          <div className="comments-section">
            <h2>Comments</h2>
            <div className="comments-list">
              {post.comments.map((comment) => (
                <div key={comment.id} className="comment">
                  <div className="comment-avatar">
                    {comment.user?.username?.[0]?.toUpperCase() || 'U'}
                  </div>
                  <div className="comment-content">
                    <div className="comment-header">
                      <span className="comment-username">
                        {comment.user?.username || 'Unknown'}
                      </span>
                      <span className="comment-date">
                        {formatDate(comment.created_at)}
                      </span>
                    </div>
                    <p>{comment.content}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </article>
    </div>
  );
}
