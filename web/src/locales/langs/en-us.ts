const local: App.I18n.Schema = {
  system: {
    title: 'ModelGate',
    updateTitle: 'System Version Update Notification',
    updateContent: 'A new version of the system has been detected. Do you want to refresh the page immediately?',
    updateConfirm: 'Refresh immediately',
    updateCancel: 'Later'
  },
  common: {
    action: 'Action',
    add: 'Add',
    requestFailed: 'Request Failed',
    addSuccess: 'Add Success',
    addFailed: 'Add Failed',
    sendSuccess: 'Send Success',
    copySuccess: 'Copy Success',
    sendFailed: 'Send Failed',
    backToHome: 'Back to home',
    batchDelete: 'Batch Delete',
    cancel: 'Cancel',
    close: 'Close',
    check: 'Check',
    expandColumn: 'Expand Column',
    columnSetting: 'Column Setting',
    config: 'Config',
    confirm: 'Confirm',
    delete: 'Delete',
    deleteSuccess: 'Delete Success',
    deleteFailed: 'Delete Failed',
    confirmDelete: 'Are you sure you want to delete?',
    edit: 'Edit',
    warning: 'Warning',
    error: 'Error',
    index: 'Index',
    keywordSearch: 'Please enter keyword',
    logout: 'Logout',
    logoutConfirm: 'Are you sure you want to log out?',
    lookForward: 'Coming soon',
    modify: 'Modify',
    modifySuccess: 'Modify Success',
    modifyFailed: 'Modify Failed',
    noData: 'No Data',
    operate: 'Operate',
    optional: 'Optional',
    pleaseCheckValue: 'Please check whether the value is valid',
    refresh: 'Refresh',
    reset: 'Reset',
    search: 'Search',
    switch: 'Switch',
    tip: 'Tip',
    trigger: 'Trigger',
    update: 'Update',
    updateSuccess: 'Update Success',
    updateFailed: 'Update Failed',
    userCenter: 'User Center',
    yesOrNo: {
      yes: 'Yes',
      no: 'No'
    },
    status: {
      enabled: 'Enabled',
      disabled: 'Disabled'
    },
    apiMethod: {
      get: 'GET',
      post: 'POST',
      put: 'PUT',
      delete: 'DELETE',
      patch: 'PATCH'
    }
  },
  request: {
    logout: 'Logout user after request failed',
    logoutMsg: 'User status is invalid, please log in again',
    logoutWithModal: 'Pop up modal after request failed and then log out user',
    logoutWithModalMsg: 'User status is invalid, please log in again',
    refreshToken: 'The requested token has expired, refresh the token',
    tokenExpired: 'The requested token has expired'
  },
  theme: {
    themeSchema: {
      title: 'Theme Schema',
      light: 'Light',
      dark: 'Dark',
      auto: 'Follow System'
    },
    grayscale: 'Grayscale',
    colourWeakness: 'Colour Weakness',
    layoutMode: {
      title: 'Layout Mode',
      vertical: 'Vertical Menu Mode',
      horizontal: 'Horizontal Menu Mode',
      'vertical-mix': 'Vertical Mix Menu Mode',
      'horizontal-mix': 'Horizontal Mix menu Mode',
      reverseHorizontalMix: 'Reverse first level menus and child level menus position'
    },
    recommendColor: 'Apply Recommended Color Algorithm',
    recommendColorDesc: 'The recommended color algorithm refers to',
    themeColor: {
      title: 'Theme Color',
      primary: 'Primary',
      info: 'Info',
      success: 'Success',
      warning: 'Warning',
      error: 'Error',
      followPrimary: 'Follow Primary'
    },
    scrollMode: {
      title: 'Scroll Mode',
      wrapper: 'Wrapper',
      content: 'Content'
    },
    page: {
      animate: 'Page Animate',
      mode: {
        title: 'Page Animate Mode',
        fade: 'Fade',
        'fade-slide': 'Slide',
        'fade-bottom': 'Fade Zoom',
        'fade-scale': 'Fade Scale',
        'zoom-fade': 'Zoom Fade',
        'zoom-out': 'Zoom Out',
        none: 'None'
      }
    },
    fixedHeaderAndTab: 'Fixed Header And Tab',
    header: {
      height: 'Header Height',
      breadcrumb: {
        visible: 'Breadcrumb Visible',
        showIcon: 'Breadcrumb Icon Visible'
      },
      multilingual: {
        visible: 'Display multilingual button'
      }
    },
    tab: {
      visible: 'Tab Visible',
      cache: 'Tag Bar Info Cache',
      height: 'Tab Height',
      mode: {
        title: 'Tab Mode',
        chrome: 'Chrome',
        button: 'Button'
      }
    },
    sider: {
      inverted: 'Dark Sider',
      width: 'Sider Width',
      collapsedWidth: 'Sider Collapsed Width',
      mixWidth: 'Mix Sider Width',
      mixCollapsedWidth: 'Mix Sider Collapse Width',
      mixChildMenuWidth: 'Mix Child Menu Width'
    },
    footer: {
      visible: 'Footer Visible',
      fixed: 'Fixed Footer',
      height: 'Footer Height',
      right: 'Right Footer'
    },
    watermark: {
      visible: 'Watermark Full Screen Visible',
      text: 'Watermark Text'
    },
    themeDrawerTitle: 'Theme Configuration',
    pageFunTitle: 'Page Function',
    resetCacheStrategy: {
      title: 'Reset Cache Strategy',
      close: 'Close Page',
      refresh: 'Refresh Page'
    },
    configOperation: {
      copyConfig: 'Copy Config',
      copySuccessMsg: 'Copy Success, Please replace the variable "themeSettings" in "src/theme/settings.ts"',
      resetConfig: 'Reset Config',
      resetSuccessMsg: 'Reset Success'
    }
  },
  route: {
    login: 'Login',
    403: 'No Permission',
    404: 'Page Not Found',
    500: 'Server Error',
    'iframe-page': 'Iframe',
    home: 'Home',
    document: 'Document',
    document_project: 'Project Document',
    'document_project-link': 'Project Document(External Link)',
    document_vue: 'Vue Document',
    document_vite: 'Vite Document',
    document_unocss: 'UnoCSS Document',
    document_naive: 'Naive UI Document',
    document_antd: 'Ant Design Vue Document',
    document_alova: 'Alova Document',
    'user-center': 'User Center',
    about: 'About',
    function: 'System Function',
    alova: 'Alova Example',
    alova_request: 'Alova Request',
    alova_user: 'User List',
    alova_scenes: 'Scenario Request',
    function_tab: 'Tab',
    'function_multi-tab': 'Multi Tab',
    'function_hide-child': 'Hide Child',
    'function_hide-child_one': 'Hide Child',
    'function_hide-child_two': 'Two',
    'function_hide-child_three': 'Three',
    function_request: 'Request',
    'function_toggle-auth': 'Toggle Auth',
    'function_super-page': 'Super Admin Visible',
    manage: 'System Manage',
    manage_user: 'User Manage',
    'manage_user-detail': 'User Detail',
    manage_role: 'Role Manage',
    manage_menu: 'Menu Manage',
    'multi-menu': 'Multi Menu',
    'multi-menu_first': 'Menu One',
    'multi-menu_first_child': 'Menu One Child',
    'multi-menu_second': 'Menu Two',
    'multi-menu_second_child': 'Menu Two Child',
    'multi-menu_second_child_home': 'Menu Two Child Home',
    exception: 'Exception',
    exception_403: '403',
    exception_404: '404',
    exception_500: '500',
    plugin: 'Plugin',
    plugin_copy: 'Copy',
    plugin_charts: 'Charts',
    plugin_charts_echarts: 'ECharts',
    plugin_charts_antv: 'AntV',
    plugin_charts_vchart: 'VChart',
    plugin_editor: 'Editor',
    plugin_editor_quill: 'Quill',
    plugin_editor_markdown: 'Markdown',
    plugin_icon: 'Icon',
    plugin_map: 'Map',
    plugin_print: 'Print',
    plugin_swiper: 'Swiper',
    plugin_video: 'Video',
    plugin_barcode: 'Barcode',
    plugin_pinyin: 'pinyin',
    plugin_excel: 'Excel',
    plugin_pdf: 'PDF preview',
    plugin_gantt: 'Gantt Chart',
    plugin_gantt_dhtmlx: 'dhtmlxGantt',
    plugin_gantt_vtable: 'VTableGantt',
    plugin_typeit: 'Typeit',
    plugin_tables: 'Tables',
    plugin_tables_vtable: 'VTable',
    relay: 'Model Management',
    relay_provider: 'Provider',
    relay_model: 'Model List',
    "relay_model-pricing": 'Model Pricing',
    "relay_provider-api-key": 'Provider API Key',
    "user": 'User',
    "user_account": 'Account',
    "user_api-key": 'User API Key',
    "usage": 'Usage',
    "usage_request": 'Request',
    "usage_ledger": 'Ledger',
  },
  page: {
    login: {
      common: {
        loginOrRegister: 'Login / Register',
        usernamePlaceholder: 'Please enter user name',
        phonePlaceholder: 'Please enter phone number',
        codePlaceholder: 'Please enter verification code',
        passwordPlaceholder: 'Please enter password',
        confirmPasswordPlaceholder: 'Please enter password again',
        codeLogin: 'Verification code login',
        confirm: 'Confirm',
        back: 'Back',
        validateSuccess: 'Verification passed',
        loginSuccess: 'Login successfully',
        loginFailed: 'Login failed',
        welcomeBack: 'Welcome back, {userName} !'
      },
      pwdLogin: {
        title: 'Password Login',
        rememberMe: 'Remember me',
        forgetPassword: 'Forget password?',
        register: 'Register',
        otherAccountLogin: 'Other Account Login',
        otherLoginMode: 'Other Login Mode',
        superAdmin: 'Super Admin',
        admin: 'Admin',
        user: 'User'
      },
      codeLogin: {
        title: 'Verification Code Login',
        getCode: 'Get verification code',
        reGetCode: 'Reacquire after {time}s',
        sendCodeSuccess: 'Verification code sent successfully',
        imageCodePlaceholder: 'Please enter image verification code'
      },
      register: {
        title: 'Register',
        agreement: 'I have read and agree to',
        protocol: '《User Agreement》',
        policy: '《Privacy Policy》'
      },
      resetPwd: {
        title: 'Reset Password'
      },
      bindWeChat: {
        title: 'Bind WeChat'
      }
    },
    about: {
      title: 'About',
      introduction: `SoybeanAdmin is an elegant and powerful admin template, based on the latest front-end technology stack, including Vue3, Vite5, TypeScript, Pinia and UnoCSS. It has built-in rich theme configuration and components, strict code specifications, and an automated file routing system. In addition, it also uses the online mock data solution based on ApiFox. SoybeanAdmin provides you with a one-stop admin solution, no additional configuration, and out of the box. It is also a best practice for learning cutting-edge technologies quickly.`,
      projectInfo: {
        title: 'Project Info',
        version: 'Version',
        latestBuildTime: 'Latest Build Time',
        githubLink: 'Github Link',
        previewLink: 'Preview Link'
      },
      prdDep: 'Production Dependency',
      devDep: 'Development Dependency'
    },
    home: {
      branchDesc:
        'For the convenience of everyone in developing and updating the merge, we have streamlined the code of the main branch, only retaining the homepage menu, and the rest of the content has been moved to the example branch for maintenance. The preview address displays the content of the example branch.',
      greeting: 'Good morning, {userName}, today is another day full of vitality!',
      weatherDesc: 'Today is cloudy to clear, 20℃ - 25℃!',
      providerCount: 'Provider Count',
      modelCount: 'Model Count',
      apiKeyCount: 'Provider API Key Count',
      downloadCount: 'Download Count',
      registerCount: 'Register Count',
      schedule: 'Work and rest Schedule',
      study: 'Study',
      work: 'Work',
      rest: 'Rest',
      entertainment: 'Entertainment',
      requestCount: 'Request Count',
      pointAmount: 'Point Amount',
      requestVolume: 'Request Volume',
      pointVolume: 'Point Volume',
      timeRange: 'Time Range',
      today: 'Today',
      last7Days: 'Last 7 Days',
      last30Days: 'Last 30 Days',
      custom: 'Custom',
      chartTab: {
        request: 'Request',
        point: 'Point'
      },
      projectNews: {
        title: 'Project News',
        moreNews: 'More News',
        desc1: 'just wrote some of the workbench pages casually, and it was enough to see!'
      },
      creativity: 'Creativity'
    },
    function: {
      tab: {
        tabOperate: {
          title: 'Tab Operation',
          addTab: 'Add Tab',
          addTabDesc: 'To about page',
          closeTab: 'Close Tab',
          closeCurrentTab: 'Close Current Tab',
          closeAboutTab: 'Close "About" Tab',
          addMultiTab: 'Add Multi Tab',
          addMultiTabDesc1: 'To MultiTab page',
          addMultiTabDesc2: 'To MultiTab page(with query params)'
        },
        tabTitle: {
          title: 'Tab Title',
          changeTitle: 'Change Title',
          change: 'Change',
          resetTitle: 'Reset Title',
          reset: 'Reset'
        }
      },
      multiTab: {
        routeParam: 'Route Param',
        backTab: 'Back function_tab'
      },
      toggleAuth: {
        toggleAccount: 'Toggle Account',
        authHook: 'Auth Hook Function `hasAuth`',
        superAdminVisible: 'Super Admin Visible',
        adminVisible: 'Admin Visible',
        adminOrUserVisible: 'Admin and User Visible'
      },
      request: {
        repeatedErrorOccurOnce: 'Repeated Request Error Occurs Once',
        repeatedError: 'Repeated Request Error',
        repeatedErrorMsg1: 'Custom Request Error 1',
        repeatedErrorMsg2: 'Custom Request Error 2'
      }
    },
    alova: {
      scenes: {
        captchaSend: 'Captcha Send',
        autoRequest: 'Auto Request',
        visibilityRequestTips: 'Automatically request when switching browser window',
        pollingRequestTips: 'It will request every 3 seconds',
        networkRequestTips: 'Automatically request after network reconnecting',
        refreshTime: 'Refresh Time',
        startRequest: 'Start Request',
        stopRequest: 'Stop Request',
        requestCrossComponent: 'Request Cross Component',
        triggerAllRequest: 'Manually Trigger All Automated Requests'
      }
    },
    manage: {
      common: {
        status: {
          enable: 'Enable',
          disable: 'Disable'
        }
      },
      permission: {
        title: 'Permission List',
        name: 'Permission Name',
        code: 'Permission Code',
        data: 'Permission Data',
        desc: 'Permission Description',
        form: {
          name: 'Please enter permission name',
          code: 'Please enter permission code',
          data: 'Please enter permission data',
          desc: 'Please enter permission description',
          path: 'Please enter path',
          method: 'Please select method'
        },
        addPermission: 'Add Permission',
        editPermission: 'Edit Permission',
        addApiPerm: 'Add API Permission'
      },
      role: {
        title: 'Role List',
        name: 'Role Name',
        code: 'Role Code',
        status: 'Role Status',
        description: 'Role Description',
        menuAuth: 'Menu Auth',
        buttonAuth: 'Button Auth',
        isSuperAdmin: 'Is Super Admin',
        form: {
          name: 'Please enter role name',
          code: 'Please enter role code',
          status: 'Please select role status',
          description: 'Please enter role description',
          isSuperAdmin: 'Please select is super admin'
        },
        addRole: 'Add Role',
        editRole: 'Edit Role'
      },
      user: {
        title: 'User List',
        username: 'User Name',
        gender: 'Gender',
        nickname: 'Nick Name',
        phone: 'Phone Number',
        email: 'Email',
        password: 'Password',
        status: 'User Status',
        role: 'User Role',
        form: {
          username: 'Please enter user name',
          gender: 'Please select gender',
          nickname: 'Please enter nick name',
          phone: 'Please enter phone number',
          email: 'Please enter email',
          password: 'Please enter password',
          status: 'Please select user status',
          role: 'Please select user role'
        },
        addUser: 'Add User',
        editUser: 'Edit User',
        userGender: {
          male: 'Male',
          female: 'Female',
          unknown: 'Unknown'
        }
      },
      menu: {
        home: 'Home',
        title: 'Menu List',
        id: 'ID',
        parentId: 'Parent ID',
        menuType: 'Menu Type',
        menuName: 'Menu Name',
        routeName: 'Route Name',
        routePath: 'Route Path',
        pathParam: 'Path Param',
        layout: 'Layout Component',
        page: 'Page Component',
        i18nKey: 'I18n Key',
        icon: 'Icon',
        localIcon: 'Local Icon',
        iconTypeTitle: 'Icon Type',
        order: 'Order',
        constant: 'Constant',
        keepAlive: 'Keep Alive',
        href: 'Href',
        hideInMenu: 'Hide In Menu',
        activeMenu: 'Active Menu',
        multiTab: 'Multi Tab',
        fixedIndexInTab: 'Fixed Index In Tab',
        query: 'Query Params',
        button: 'Button',
        buttonCode: 'Button Code',
        buttonDesc: 'Button Desc',
        menuStatus: 'Menu Status',
        form: {
          home: 'Please select home',
          menuType: 'Please select menu type',
          menuName: 'Please enter menu name',
          routeName: 'Please enter route name',
          routePath: 'Please enter route path',
          pathParam: 'Please enter path param',
          page: 'Please select page component',
          layout: 'Please select layout component',
          i18nKey: 'Please enter i18n key',
          icon: 'Please enter iconify name',
          localIcon: 'Please enter local icon name',
          order: 'Please enter order',
          keepAlive: 'Please select whether to cache route',
          href: 'Please enter href',
          hideInMenu: 'Please select whether to hide menu',
          activeMenu: 'Please select route name of the highlighted menu',
          multiTab: 'Please select whether to support multiple tabs',
          fixedInTab: 'Please select whether to fix in the tab',
          fixedIndexInTab: 'Please enter the index fixed in the tab',
          queryKey: 'Please enter route parameter Key',
          queryValue: 'Please enter route parameter Value',
          button: 'Please select whether it is a button',
          buttonCode: 'Please enter button code',
          buttonDesc: 'Please enter button description',
          menuStatus: 'Please select menu status',
          selectPermission: 'Please select permission',
        },
        addMenu: 'Add Menu',
        editMenu: 'Edit Menu',
        addChildMenu: 'Add Child Menu',
        type: {
          unspecified: 'Unspecified',
          directory: 'Directory',
          menu: 'Menu'
        },
        iconType: {
          iconify: 'Iconify Icon',
          local: 'Local Icon'
        }
      }
    },
    relay: {
      common: {
        currency: {
          usd: 'USD',
          cny: 'CNY',
          point: 'Point'
        },
        enableStatus: {
          enabled: 'Enabled',
          disabled: 'Disabled'
        },
        provider: {
          status: {
            enabled: 'Enabled',
            disabled: 'Disabled'
          }
        },
        apiKey: {
          status: {
            enabled: 'Enabled',
            disabled: 'Disabled',
            cooldown: 'Cooldown',
            revoked: 'Revoked'
          }
        },
        model: {
          status: {
            enabled: 'Enabled',
            disabled: 'Disabled',
            deprecated: 'Deprecated'
          }
        }
      },
      provider: {
        addProvider: 'Add Provider',
        editProvider: 'Edit Provider',
        title: 'Provider',
        name: 'Name',
        code: 'Code',
        baseUrl: 'Base URL',
        status: 'Status',
        form: {
          name: 'Name',
          code: 'Code',
          baseUrl: 'Base URL',
          status: 'Status',
        }
      },
      providerApiKey: {
        title: 'Provider API Key',
        providerId: 'Provider ID',
        providerCode: 'Provider Code',
        name: 'Name',
        key: 'Key',
        weight: 'Weight',
        status: 'Status',
        lastUsedAt: 'Last Used At',
        form: {
          providerId: 'Provider ID',
          providerCode: 'Provider Code',
          name: 'Name',
          key: 'Key',
          weight: 'Weight',
          status: 'Status',
        }
      },
      model: {
        title: 'Model',
        providerId: 'Provider ID',
        providerCode: 'Provider Code',
        code: 'Code',
        actualCode: 'Actual Code',
        name: 'Name',
        priority: 'Priority',
        weight: 'Weight',
        status: 'Status',
        form: {
          providerId: 'Provider ID',
          code: 'Code',
          actualCode: 'Actual Code',
          providerCode: 'Provider Code',
          name: 'Name',
          priority: 'Priority',
          weight: 'Weight',
          status: 'Status',
        }
      },
      modelPricing: {
        title: 'Model Pricing',
        providerCode: 'Provider Code',
        modelCode: 'Model Code',
        currency: 'Currency',
        pointsPerCurrency: 'Points Per Currency',
        tokenNum: 'Token Num',
        inputPrice: 'Input Price',
        inputCachePrice: 'Input Cache Price',
        outputPrice: 'Output Price',
        effectiveFrom: 'Effective From',
        effectiveTo: 'Effective To',
        status: 'Status',
        form: {
          providerCode: 'Provider Code',
          modelCode: 'Model Code',
          currency: 'Currency',
          pointsPerCurrency: 'Points Per Currency',
          tokenNum: 'Token Num',
          inputPrice: 'Input Price',
          inputCachePrice: 'Input Cache Price',
          outputPrice: 'Output Price',
          effectiveFrom: 'Effective From',
          effectiveTo: 'Effective To',
          status: 'Status',
        }
      },

    },
    usage: {
      common: {
        request: {
          status: {
            pending: 'Pending',
            success: 'Success',
            failed: 'Failed',
            cancelled: 'Cancelled',
          }
        },
        ledger: {
          type: {
            consume: 'Consume',
            refund: 'Refund',
            charge: 'Charge',
            adjust: 'Adjust',
          }
        }
      },
      request: {
        title: 'Request List',
        requestUuid: 'Request UUID',
        accountId: 'Account ID',
        accountName: 'Account Name',
        providerCode: 'Provider Code',
        promptTokens: 'Prompt Tokens',
        completionTokens: 'Completion Tokens',
        totalTokens: 'Total Tokens',
        modelCode: 'Model Code',
        status: 'Status',
        createdAt: 'Request Time',
        elapsedTime: 'Elapsed Time(seconds)',
        form: {
          providerCode: 'Provider Code',
          accountId: 'Account',
          status: 'Status',
        }
      },
      ledger: {
        title: 'Ledger List',
        accountId: 'Account ID',
        accountName: 'Account Name',
        type: 'Type',
        amount: 'Amount',
        balanceAfter: 'Balance After',
        requestId: 'Request ID',
        createdAt: 'Created At',
        form: {
          accountId: 'Account',
          type: 'Type',
          code: 'Code'
        }
      },
    },
    user: {
      account: {
        title: 'Account List',
        name: 'Name',
        nickname: 'Nickname',
        balance: 'Balance',
        status: 'Status',
        createdAt: 'Created At',
        updatedAt: 'Updated At',
        addAccount: 'Add Account',
        editAccount: 'Edit Account',
        form: {
          name: 'Name',
          nickname: 'Nickname',
          balance: 'Balance',
          status: 'Status',
        }
      },
      apiKey: {
        title: 'API Key List',
        addApiKey: 'Add API Key',
        editApiKey: 'Edit API Key',
        accountName: 'Account Name',
        keyName: 'Key Name',
        key: 'Key',
        scope: 'Scope',
        quoteLimit: 'Quote Limit',
        quoteUsed: 'Quote Used',
        rateLimit: 'Rate Limit',
        lastUsedAt: 'Last Used At',
        expiredAt: 'Expired At',
        status: 'Status',
        remark: 'Remark',
        form: {
          accountId: 'Please select account',
          keyName: 'Please enter key name',
          scope: 'Please select scope',
          quoteLimit: 'Please enter quote limit',
          rateLimit: 'Please enter rate limit',
          expiredAt: 'Please select expired at',
          status: 'Please select status',
          remark: 'Please enter remark',
        },
        keyModal: {
          title: 'API Key Generated Successfully',
          warningTitle: 'Important Notice',
          warningMessage: 'This API Key will only be displayed once. Please copy and save it immediately!',
          yourKey: 'Your API Key:',
          copied: 'Copied',
          copy: 'Copy',
          confirm: 'I have copied, close'
        }
      }
    },
    userCenter: {
      title: 'User Center'
    },
  },
  form: {
    required: 'Cannot be empty',
    username: {
      required: 'Please enter user name',
      invalid: 'User name format is incorrect'
    },
    phone: {
      required: 'Please enter phone number',
      invalid: 'Phone number format is incorrect'
    },
    pwd: {
      required: 'Please enter password',
      invalid: '6-18 characters, including letters, numbers, and underscores'
    },
    confirmPwd: {
      required: 'Please enter password again',
      invalid: 'The two passwords are inconsistent'
    },
    code: {
      required: 'Please enter verification code',
      invalid: 'Verification code format is incorrect'
    },
    email: {
      required: 'Please enter email',
      invalid: 'Email format is incorrect'
    }
  },
  dropdown: {
    closeCurrent: 'Close Current',
    closeOther: 'Close Other',
    closeLeft: 'Close Left',
    closeRight: 'Close Right',
    closeAll: 'Close All'
  },
  icon: {
    themeConfig: 'Theme Configuration',
    themeSchema: 'Theme Schema',
    lang: 'Switch Language',
    fullscreen: 'Fullscreen',
    fullscreenExit: 'Exit Fullscreen',
    reload: 'Reload Page',
    collapse: 'Collapse Menu',
    expand: 'Expand Menu',
    pin: 'Pin',
    unpin: 'Unpin'
  },
  datatable: {
    itemCount: 'Total {total} items'
  }
};

export default local;
