import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import App from './App.tsx';
import './index.css';
import { AuthProvider } from './lib/auth.tsx';
import { ProtectedRoute } from './components/ProtectedRoute.tsx';
import { Home } from './pages/Home.tsx';
import { Login } from './pages/Login.tsx';
import { Register } from './pages/Register.tsx';
import { Feed } from './pages/Feed.tsx';
import { Explore } from './pages/Explore.tsx';
import { Profile } from './pages/Profile.tsx';
import { CreatePost } from './pages/CreatePost.tsx';
import { PostDetail } from './pages/PostDetail.tsx';
import { ConfirmUser } from './pages/ConfirmUser.tsx';

export const BASE_API = 'http://localhost:8000/v1';

const router = createBrowserRouter([
  {
    path: '/',
    element: <App />,
    children: [
      {
        index: true,
        element: <Home />,
      },
      {
        path: 'login',
        element: <Login />,
      },
      {
        path: 'register',
        element: <Register />,
      },
      {
        path: 'confirm/:token',
        element: <ConfirmUser />,
      },
      {
        path: 'feed',
        element: (
          <ProtectedRoute>
            <Feed />
          </ProtectedRoute>
        ),
      },
      {
        path: 'explore',
        element: (
          <ProtectedRoute>
            <Explore />
          </ProtectedRoute>
        ),
      },
      {
        path: 'profile/:userId',
        element: (
          <ProtectedRoute>
            <Profile />
          </ProtectedRoute>
        ),
      },
      {
        path: 'create-post',
        element: (
          <ProtectedRoute>
            <CreatePost />
          </ProtectedRoute>
        ),
      },
      {
        path: 'posts/:postId',
        element: (
          <ProtectedRoute>
            <PostDetail />
          </ProtectedRoute>
        ),
      },
    ],
  },
]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AuthProvider>
      <RouterProvider router={router} />
    </AuthProvider>
  </StrictMode>
);
