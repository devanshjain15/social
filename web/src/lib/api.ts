import { BASE_API } from '../main';

export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
  is_active: boolean;
  is_following?: boolean; // Added for follow status
}

export interface Post {
  id: number;
  title: string;
  content: string;
  tags: string[];
  user_id: number;
  created_at: string;
  updated_at: string;
  version: number;
  user?: User;
  comments?: Comment[];
  comments_count?: number;
}

export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  content: string;
  created_at: string;
  user?: User;
}

export interface FeedResponse {
  posts: Post[];
  pagination: {
    limit: number;
    offset: number;
    total: number;
  };
}

class ApiClient {
  private token: string | null = null;

  constructor() {
    this.token = localStorage.getItem('auth_token');
  }

  setToken(token: string) {
    this.token = token;
    localStorage.setItem('auth_token', token);
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem('auth_token');
  }

  getToken() {
    return this.token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (options.headers) {
      Object.assign(headers, options.headers);
    }

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${BASE_API}${endpoint}`, {
      ...options,
      headers,
      credentials: 'include', // Always send cookies
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(error.error || `HTTP error! status: ${response.status}`);
    }

    if (response.status === 204) {
      return {} as T;
    }

    const json = await response.json();

    // Backend wraps responses in {"data": ...}, unwrap it
    if (json && typeof json === 'object' && 'data' in json) {
      return json.data as T;
    }

    return json as T;
  }

  // Auth endpoints
  async register(username: string, email: string, password: string) {
    return this.request<{
      message: string;
      activation_url: string;
      activation_token: string;
    }>('/authentication/user', {
      method: 'POST',
      body: JSON.stringify({ username, email, password }),
    });
  }

  async login(email: string, password: string) {
    const response = await this.request<{ token: string; user?: User }>('/authentication/token', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });

    this.setToken(response.token);
    return response.user;
  }

  async logout() {
    return this.request<void>('/authentication/logout', {
      method: 'POST',
    });
  }

  async activateUser(token: string) {
    return this.request<void>(`/users/activate/${token}`, {
      method: 'PUT',
    });
  }

  // User endpoints
  async getAllUsers() {
    return this.request<User[]>('/users');
  }

  async getUser(userId: number) {
    return this.request<User>(`/users/${userId}`);
  }

  async followUser(userId: number) {
    return this.request<void>(`/users/${userId}/follow`, {
      method: 'PUT',
    });
  }

  async unfollowUser(userId: number) {
    return this.request<void>(`/users/${userId}/unfollow`, {
      method: 'PUT',
    });
  }

  async getUserFeed(limit = 10, offset = 0, sort = 'desc') {
    // Backend returns array directly after unwrapping
    return this.request<Post[]>(`/users/feed?limit=${limit}&offset=${offset}&sort=${sort}`);
  }

  // Post endpoints
  async createPost(title: string, content: string, tags: string[] = []) {
    return this.request<Post>('/posts/', {
      method: 'POST',
      body: JSON.stringify({ title, content, tags }),
    });
  }

  async getPost(postId: number) {
    return this.request<Post>(`/posts/${postId}`);
  }

  async updatePost(postId: number, title?: string, content?: string) {
    const body: { title?: string; content?: string } = {};
    if (title !== undefined) body.title = title;
    if (content !== undefined) body.content = content;

    return this.request<Post>(`/posts/${postId}`, {
      method: 'PATCH',
      body: JSON.stringify(body),
    });
  }

  async deletePost(postId: number) {
    return this.request<void>(`/posts/${postId}`, {
      method: 'DELETE',
    });
  }

  // Comment endpoints
  async createComment(postId: number, content: string) {
    return this.request<Comment>(`/posts/${postId}/comments`, {
      method: 'POST',
      body: JSON.stringify({ content }),
    });
  }
}

export const api = new ApiClient();
