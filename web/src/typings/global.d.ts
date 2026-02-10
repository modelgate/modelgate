export {};

declare global {
  export interface Window {
    /** NProgress instance */
    NProgress?: import('nprogress').NProgress;
    /** Loading bar instance */
    $loadingBar?: import('naive-ui').LoadingBarProviderInst;
    /** Dialog instance */
    $dialog?: import('naive-ui').DialogProviderInst;
    /** Message instance */
    $message?: import('naive-ui').MessageProviderInst;
    /** Notification instance */
    $notification?: import('naive-ui').NotificationProviderInst;
    /** Build time of the project */
    BUILD_TIME?: string;
    /** App version from package.json */
    APP_VERSION?: string;
  }

  /** Build time of the project */
  export const BUILD_TIME: string;
  /** App version from package.json */
  export const APP_VERSION: string;
}
