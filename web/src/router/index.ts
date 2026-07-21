import { type RouteRecordRaw, createRouter } from "vue-router"
import { history, flatMultiLevelRoutes } from "./helper"
import routeSettings from "@/config/route"

const Layouts = () => import("@/layouts/index.vue")

/**
 * 常驻路由
 * 除了 redirect/403/404/login 等隐藏页面，其他页面建议设置 Name 属性
 */
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: "/",
    redirect: "/board/summary",
    meta: {
      hidden: true
    }
  },
  {
    path: "/redirect",
    component: Layouts,
    meta: {
      hidden: true
    },
    children: [
      {
        path: ":path(.*)",
        component: () => import("@/views/redirect/index.vue")
      }
    ]
  },
  {
    path: "/403",
    component: () => import("@/views/error-page/403.vue"),
    meta: {
      hidden: true
    }
  },
  {
    path: "/404",
    component: () => import("@/views/error-page/404.vue"),
    meta: {
      hidden: true
    },
    alias: "/:pathMatch(.*)*"
  },
  {
    path: "/login",
    component: () => import("@/views/login/index.vue"),
    meta: {
      hidden: true
    },
    name: "Login"
  },
  {
    path: "/oauth/callback/:provider?",
    component: () => import("@/views/login/oauth-callback.vue"),
    meta: {
      hidden: true
    },
    name: "OAuthCallback"
  },
  {
    path: "/profile",
    component: () => import("@/views/profile/index.vue"),
    meta: {
      hidden: true
    },
    name: "Profile"
  },

  {
    path: "/navi",
    component: () => import("@/views/navi/index.vue"),
    name: "Navigation",
    meta: {
      hidden: true
    }
  },
  {
    path: "/board",
    component: Layouts,
    redirect: "/board/summary",
    meta: {
      title: "看板",
      svgIcon: "dashboard"
    },
    children: [
      {
        path: "summary",
        component: () => import("@/views/cluster/index.vue"),
        name: "Summary",
        meta: {
          title: "集群概览",
          affix: false
        }
      },

      {
        path: "minikube",
        component: () => import("@/views/minikube/index.vue"),
        name: "minikube",
        meta: {
          title: "集群安装"
        }
      },
      {
        path: "navy",
        component: () => import("@/views/navy/index.vue"),
        name: "Navy",
        meta: {
          title: "运维导航"
        }
      }
    ]
  },

  {
    path: "/link",
    component: Layouts,
    meta: {
      title: "项目管理",

      svgIcon: "docs"
    },
    children: [
      {
        path: "https://cilikube.cillian.website",
        component: () => { },
        name: "Link1",
        meta: {
          title: "中文文档"
        }
      },
      {
        path: "https://www.cillian.website",
        component: () => { },
        name: "Link2",
        meta: {
          title: "我的博客"
        }
      },
      {
        path: "techstack",
        component: () => import("@/views/techstack/index.vue"),
        name: "技术栈",
        meta: {
          title: "技术栈"
        }
      }
    ]
  },

  // {
  //   path: "/table",
  //   component: Layouts,
  //   redirect: "/table/element-plus",
  //   name: "Table",
  //   meta: {
  //     title: "表格",
  //     elIcon: "Grid"
  //   },
  //   children: [
  //     {
  //       path: "element-plus",
  //       component: () => import("@/views/table/element-plus/index.vue"),
  //       name: "ElementPlus",
  //       meta: {
  //         title: "Element Plus",
  //         keepAlive: true
  //       }
  //     },
  //     {
  //       path: "vxe-table",
  //       component: () => import("@/views/table/vxe-table/index.vue"),
  //       name: "VxeTable",
  //       meta: {
  //         title: "Vxe Table",
  //         keepAlive: true
  //       }
  //     }
  //   ]
  // },
  // {
  //   path: "/menu",
  //   component: Layouts,
  //   redirect: "/menu/menu1",
  //   name: "Menu",
  //   meta: {
  //     title: "多级路由",
  //     svgIcon: "menu"
  //   },
  //   children: [
  //     {
  //       path: "menu1",
  //       component: () => import("@/views/menu/menu1/index.vue"),
  //       redirect: "/menu/menu1/menu1-1",
  //       name: "Menu1",
  //       meta: {
  //         title: "menu1"
  //       },
  //       children: [
  //         {
  //           path: "menu1-1",
  //           component: () => import("@/views/menu/menu1/menu1-1/index.vue"),
  //           name: "Menu1-1",
  //           meta: {
  //             title: "menu1-1",
  //             keepAlive: true
  //           }
  //         },
  //         {
  //           path: "menu1-2",
  //           component: () => import("@/views/menu/menu1/menu1-2/index.vue"),
  //           redirect: "/menu/menu1/menu1-2/menu1-2-1",
  //           name: "Menu1-2",
  //           meta: {
  //             title: "menu1-2"
  //           },
  //           children: [
  //             {
  //               path: "menu1-2-1",
  //               component: () => import("@/views/menu/menu1/menu1-2/menu1-2-1/index.vue"),
  //               name: "Menu1-2-1",
  //               meta: {
  //                 title: "menu1-2-1",
  //                 keepAlive: true
  //               }
  //             },
  //             {
  //               path: "menu1-2-2",
  //               component: () => import("@/views/menu/menu1/menu1-2/menu1-2-2/index.vue"),
  //               name: "Menu1-2-2",
  //               meta: {
  //                 title: "menu1-2-2",
  //                 keepAlive: true
  //               }
  //             }
  //           ]
  //         },
  //         {
  //           path: "menu1-3",
  //           component: () => import("@/views/menu/menu1/menu1-3/index.vue"),
  //           name: "Menu1-3",
  //           meta: {
  //             title: "menu1-3",
  //             keepAlive: true
  //           }
  //         }
  //       ]
  //     },
  //     {
  //       path: "menu2",
  //       component: () => import("@/views/menu/menu2/index.vue"),
  //       name: "Menu2",
  //       meta: {
  //         title: "menu2",
  //         keepAlive: true
  //       }
  //     }
  //   ]
  // },
  // {
  //   path: "/hook-demo",
  //   component: Layouts,
  //   redirect: "/hook-demo/use-fetch-select",
  //   name: "HookDemo",
  //   meta: {
  //     title: "Hook",
  //     elIcon: "Menu",
  //     alwaysShow: true
  //   },
  //   children: [
  //     {
  //       path: "use-fetch-select",
  //       component: () => import("@/views/hook-demo/use-fetch-select.vue"),
  //       name: "UseFetchSelect",
  //       meta: {
  //         title: "useFetchSelect"
  //       }
  //     },
  //     {
  //       path: "use-fullscreen-loading",
  //       component: () => import("@/views/hook-demo/use-fullscreen-loading.vue"),
  //       name: "UseFullscreenLoading",
  //       meta: {
  //         title: "useFullscreenLoading"
  //       }
  //     },
  //     {
  //       path: "use-watermark",
  //       component: () => import("@/views/hook-demo/use-watermark.vue"),
  //       name: "UseWatermark",
  //       meta: {
  //         title: "useWatermark"
  //       }
  //     }
  //   ]
  // }
]

/**
 * 动态路由
 * 用来放置有权限 (Roles 属性) 的路由
 * 必须带有 Name 属性
 */
export const dynamicRoutes: RouteRecordRaw[] = [
  // Kubernetes 集群管理路由
  {
    path: "/cluster",
    component: Layouts,
    meta: {
      title: "集群",
      svgIcon: "cluster",
      roles: ["admin", "editor", "viewer"] // 所有角色都可以访问集群信息
    },
    children: [
      {
        path: "management",
        component: () => import("@/views/cluster/management/index.vue"),
        name: "ClusterManagement",
        meta: {
          title: "集群管理",
          roles: ["admin"] // 只有管理员可以管理集群
        }
      },
      {
        path: "node",
        component: () => import("@/views/node/index.vue"),
        name: "Node",
        meta: {
          title: "节点",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看节点
        }
      },
      {
        path: "namespace",
        component: () => import("@/views/namespace/index.vue"),
        name: "Namespace",
        meta: {
          title: "命名空间",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看命名空间
        }
      },
      {
        path: "events",
        component: () => import("@/views/events/index.vue"),
        name: "ClusterEvents",
        meta: {
          title: "集群事件",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看事件
        }
      },
      {
        path: "crd",
        component: () => import("@/views/crd/index.vue"),
        name: "CRD",
        meta: {
          title: "自定义资源",
          roles: ["admin", "editor"] // 管理员和编辑者可以查看CRD
        }
      }
    ]
  },
  // Kubernetes 工作负载路由
  {
    path: "/workloads",
    component: Layouts,
    meta: {
      title: "工作负载",
      svgIcon: "kubernetes",
      roles: ["admin", "editor", "viewer"] // 所有角色都可以访问工作负载
    },
    children: [
      {
        path: "pods",
        component: () => import("@/views/pods/index.vue"),
        name: "Pods",
        meta: {
          title: "Pod",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看Pod
        }
      },
      {
        path: "deployments",
        component: () => import("@/views/deployments/index.vue"),
        name: "Deployments",
        meta: {
          title: "Deployment",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看Deployment
        }
      }
    ]
  },
  // Kubernetes 存储路由
  {
    path: "/storage",
    component: Layouts,
    meta: {
      title: "存储",
      svgIcon: "storage",
      roles: ["admin", "editor", "viewer"] // 所有角色都可以访问存储
    },
    children: [
      {
        path: "persistentvolume",
        component: () => import("@/views/persistentvolume/index.vue"),
        name: "PersistentVolume",
        meta: {
          title: "PV",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看PV
        }
      },
      {
        path: "persistentvolumeclaim",
        component: () => import("@/views/persistentvolumeclaim/index.vue"),
        name: "PersistentVolumeClaim",
        meta: {
          title: "PVC",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看PVC
        }
      }
    ]
  },
  // Kubernetes 网络路由
  {
    path: "/network",
    component: Layouts,
    meta: {
      title: "网络",
      svgIcon: "network",
      roles: ["admin", "editor", "viewer"] // 所有角色都可以访问网络
    },
    children: [
      {
        path: "service",
        component: () => import("@/views/service/index.vue"),
        name: "Service",
        meta: {
          title: "Service",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看Service
        }
      },
      {
        path: "ingress",
        component: () => import("@/views/ingress/index.vue"),
        name: "Ingress",
        meta: {
          title: "Ingress",
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看Ingress
        }
      }
    ]
  },
  // Kubernetes 配置管理路由
  {
    path: "/config",
    component: Layouts,
    meta: {
      title: "配置管理",
      svgIcon: "config",
      roles: ["admin", "editor", "viewer"] // 所有角色都可以访问配置
    },
    children: [
      {
        path: "configmap",
        component: () => import("@/views/configmap/index.vue"),
        name: "ConfigMap",
        meta: {
          title: "ConfigMap",
          keepAlive: true,
          roles: ["admin", "editor", "viewer"] // 所有角色都可以查看ConfigMap
        }
      },
      {
        path: "secret",
        component: () => import("@/views/secret/index.vue"),
        name: "Secret",
        meta: {
          title: "Secret",
          keepAlive: true,
          roles: ["admin", "editor"] // 管理员和编辑者可以查看Secret（包含敏感信息）
        }
      }
    ]
  },
  // 权限演示路由 - 管理员和编辑者可访问
  {
    path: "/permission",
    component: Layouts,
    redirect: "/permission/page",
    name: "Permission",
    meta: {
      title: "权限演示",
      svgIcon: "lock",
      roles: ["admin", "editor"],
      alwaysShow: true
    },
    children: [
      {
        path: "page",
        component: () => import("@/views/permission/page.vue"),
        name: "PagePermission",
        meta: {
          title: "页面级权限",
          roles: ["admin"]
        }
      },
      {
        path: "directive",
        component: () => import("@/views/permission/directive.vue"),
        name: "DirectivePermission",
        meta: {
          title: "按钮级权限",
          roles: ["admin", "editor"]
        }
      }
    ]
  },
  // 系统管理路由 - 仅管理员可访问
  {
    path: "/admin",
    component: Layouts,
    redirect: "/admin/user-management",
    name: "Admin",
    meta: {
      title: "系统管理",
      svgIcon: "lock",
      roles: ["admin"],
      alwaysShow: true
    },
    children: [
      {
        path: "user-management",
        component: () => import("@/views/admin/user-management/index.vue"),
        name: "UserManagement",
        meta: {
          title: "用户管理",
          roles: ["admin"]
        }
      },
      {
        path: "role-management",
        component: () => import("@/views/admin/role-management/index.vue"),
        name: "RoleManagement",
        meta: {
          title: "角色管理",
          roles: ["admin"]
        }
      },
      {
        path: "system-settings",
        component: () => import("@/views/admin/system-settings/index.vue"),
        name: "SystemSettings",
        meta: {
          title: "系统设置",
          roles: ["admin"]
        }
      }
    ]
  }
]

const router = createRouter({
  history,
  routes: routeSettings.thirdLevelRouteCache ? flatMultiLevelRoutes(constantRoutes) : constantRoutes
})

/** 重置路由 */
export function resetRouter() {
  // 注意：所有动态路由路由必须带有 Name 属性，否则可能会不能完全重置干净
  try {
    router.getRoutes().forEach((route) => {
      const { name, meta } = route
      if (name && meta.roles?.length) {
        router.hasRoute(name) && router.removeRoute(name)
      }
    })
  } catch {
    // 强制刷新浏览器也行，只是交互体验不是很好
    window.location.reload()
  }
}
// 抛出路由实例, 在 main.js 中引用
export default router
