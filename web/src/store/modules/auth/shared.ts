import { localStg } from '@/utils/storage';
import { authServiceClient } from '@/grpc';

/** Get token */
export function getToken() {
  return localStg.get('token') || '';
}

/** Clear auth storage */
export function clearAuthStorage() {
  localStg.remove('token');
  localStg.remove('refreshToken');
}

let refreshTokenPromise: Promise<boolean> | null = null;

export async function handleRefreshToken(): Promise<boolean> {
  const rToken = localStg.get('refreshToken') || '';
  if (!rToken) {
    return false;
  }
  if (refreshTokenPromise) {
    return await refreshTokenPromise;
  }
  refreshTokenPromise = (async () => {
    try {
      const { accessToken, refreshToken } = await authServiceClient.refreshToken({
        refreshToken: rToken
      });
      if (accessToken) {
        localStg.set('token', accessToken);
        localStg.set('refreshToken', refreshToken);
      }
      return true;
    } catch {
      clearAuthStorage();
      return false;
    } finally {
      refreshTokenPromise = null;
    }
  })();
  return await refreshTokenPromise;
}
