import { localStg } from '@/utils/storage';
import { authServiceClient } from '@/grpc';

/** Get token */
export function getToken() {
  return localStg.get('token') || '';
}

/** Clear auth storage */
export function clearAuthStorage() {
  localStg.remove('token');
}

let refreshTokenPromise: Promise<boolean> | null = null;

export async function handleRefreshToken(): Promise<boolean> {
  if (refreshTokenPromise) {
    return await refreshTokenPromise;
  }
  refreshTokenPromise = (async () => {
    try {
      const { accessToken } = await authServiceClient.refreshToken({});
      if (!accessToken) {
        clearAuthStorage();
        return false;
      }
      localStg.set('token', accessToken);
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
