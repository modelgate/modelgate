import type { Interceptor, StreamResponse, UnaryResponse } from '@connectrpc/connect';
import { Code, ConnectError, createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import type { DescMessage } from '@bufbuild/protobuf';
import { AuthService } from './typings/proto/admin/v1/auth_pb';
import { RelayService } from './typings/proto/admin/v1/relay_pb';
import { SystemService } from './typings/proto/admin/v1/system_pb';
import { getToken, handleRefreshToken } from './store/modules/auth/shared';

const grpcServiceUrl = import.meta.env.VITE_PUBLIC_GRPC_SERVICE_URL || '/grpc';

const logInterceptor: Interceptor = next => async req => {
  console.log('request message: ', req.message);
  const res = await next(req);
  // console.log('response: ', res);
  return res;
};

const authInterceptor: Interceptor =
  next =>
    async (req): Promise<UnaryResponse<DescMessage, DescMessage> | StreamResponse<DescMessage, DescMessage>> => {
      const doRequest: () => Promise<
        UnaryResponse<DescMessage, DescMessage> | StreamResponse<DescMessage, DescMessage>
      > = async () => {
        const token = getToken();
        if (token) {
          req.header.set('Authorization', `Bearer ${token}`);
        }
        return await next(req);
      };
      try {
        return await doRequest();
      } catch (err) {
        if (
          err instanceof ConnectError &&
          err.code === Code.Unauthenticated &&
          !(req.method.name === 'RefreshToken' && req.method.parent.typeName === 'admin.v1.AuthService') &&
          !req.stream
        ) {
          const success = await handleRefreshToken();
          if (!success) throw err;
          return await doRequest();
        }
        console.error(err);
        throw err;
      }
    };

const transport = createConnectTransport({
  baseUrl: grpcServiceUrl,
  interceptors: [logInterceptor, authInterceptor],
  useBinaryFormat: false
  // useBinaryFormat: true
});

export const authServiceClient = createClient(AuthService, transport);
export const systemServiceClient = createClient(SystemService, transport);
export const relayServiceClient = createClient(RelayService, transport);
