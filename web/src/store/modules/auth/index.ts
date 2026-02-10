import { computed, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { defineStore } from 'pinia';
import { useLoading } from '@sa/hooks';
import { useRouterPush } from '@/hooks/common/router';
import { localStg } from '@/utils/storage';
import { authServiceClient } from '@/grpc';
import { SetupStoreId } from '@/enum';
import { $t } from '@/locales';
import type { AccessToken } from '@/typings/proto/model/system/access_token_pb';
import { useRouteStore } from '../route';
import { useTabStore } from '../tab';
import { clearAuthStorage, getToken } from './shared';

export const useAuthStore = defineStore(SetupStoreId.Auth, () => {
  const route = useRoute();
  const routeStore = useRouteStore();
  const tabStore = useTabStore();
  const { toLogin, redirectFromLogin } = useRouterPush(false);
  const { loading: loginLoading, startLoading, endLoading } = useLoading();

  const token = ref(getToken());

  const userInfo: Api.Auth.UserInfo = reactive({
    userId: '',
    userName: '',
    roles: [],
    buttons: []
  });

  /** is super role in static route */
  const isStaticSuper = computed(() => {
    const { VITE_AUTH_ROUTE_MODE, VITE_STATIC_SUPER_ROLE } = import.meta.env;

    return VITE_AUTH_ROUTE_MODE === 'static' && userInfo.roles.includes(VITE_STATIC_SUPER_ROLE);
  });

  /** Is login */
  const isLogin = computed(() => Boolean(token.value));

  /** Reset auth store */
  async function resetStore() {
    const authStore = useAuthStore();

    clearAuthStorage();

    authStore.$reset();

    if (!route.meta.constant) {
      await toLogin();
    }

    tabStore.cacheTabs();
    routeStore.resetStore();
  }

  /**
   * Login
   *
   * @param userName User name
   * @param password Password
   * @param [redirect=true] Whether to redirect after login. Default is `true`
   */
  async function login(userName: string, password: string, rememberMe: boolean, redirect = true) {
    startLoading();
    try {
      const { accessToken, refreshToken } = await authServiceClient.login({
        username: userName,
        password,
        rememberMe
      });
      if (!accessToken) {
        throw new Error('Failed to get access token');
      }
      const pass = await loginByToken(accessToken, refreshToken);
      if (pass) {
        await redirectFromLogin(redirect);

        window.$notification?.success({
          title: $t('page.login.common.loginSuccess'),
          // content: $t('page.login.common.welcomeBack', { userName: userInfo.nickname }),
          duration: 4500
        });
      }
    } catch (error: any) {
      window.$notification?.error({
        title: $t('page.login.common.loginFailed'),
        content: error.details,
        duration: 4500
      });
      resetStore();
    }

    endLoading();
  }

  async function loginByToken(accessToken: string, refreshToken: string) {
    // 1. stored in the localStorage, the later requests need it in headers
    localStg.set('token', accessToken);
    localStg.set('refreshToken', refreshToken);

    // 2. get user info
    const pass = await getUserInfo();

    if (pass) {
      token.value = accessToken;
      return true;
    }

    return false;
  }

  async function getUserInfo() {
    try {
      const { user } = await authServiceClient.getUserInfo({});
      if (!user) {
        throw new Error('Failed to get user info');
      }
      // update store
      Object.assign(userInfo, user);
      return true;
    } catch (error: any) {
      console.log(error);
      return false;
    }
  }

  async function initUserInfo() {
    const hasToken = getToken();

    if (hasToken) {
      const pass = await getUserInfo();

      if (!pass) {
        resetStore();
      }
    }
  }

  return {
    token,
    userInfo,
    isStaticSuper,
    isLogin,
    loginLoading,
    resetStore,
    login,
    initUserInfo
  };
});
