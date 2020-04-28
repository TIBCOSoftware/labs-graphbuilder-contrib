webpackJsonp(["main"],{

/***/ "./src/$$_lazy_route_resource lazy recursive":
/***/ (function(module, exports, __webpack_require__) {

var map = {
	"app/ontology.login/login.module": [
		"./src/app/ontology.login/login.module.ts",
		"login.module"
	],
	"app/ontology.project/project.module": [
		"./src/app/ontology.project/project.module.ts",
		"project.module"
	],
	"app/ontology.projects/projects.module": [
		"./src/app/ontology.projects/projects.module.ts",
		"projects.module"
	]
};
function webpackAsyncContext(req) {
	var ids = map[req];
	if(!ids)
		return Promise.reject(new Error("Cannot find module '" + req + "'."));
	return __webpack_require__.e(ids[1]).then(function() {
		return __webpack_require__(ids[0]);
	});
};
webpackAsyncContext.keys = function webpackAsyncContextKeys() {
	return Object.keys(map);
};
webpackAsyncContext.id = "./src/$$_lazy_route_resource lazy recursive";
module.exports = webpackAsyncContext;

/***/ }),

/***/ "./src/app/app-routing.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const home_layout_component_1 = __webpack_require__("./src/app/layouts/home-layout/home-layout.component.ts");
const login_layout_component_1 = __webpack_require__("./src/app/layouts/login-layout/login-layout.component.ts");
const auth_guard_service_1 = __webpack_require__("./src/app/services/auth-guard.service.ts");
exports.projectsPath = 'app/ontology.projects/projects.module#ProjectsModule';
exports.loginPath = 'app/ontology.login/login.module#LoginModule';
const appRoutes = [
    {
        path: '',
        component: home_layout_component_1.HomeLayoutComponent,
        canActivate: [auth_guard_service_1.AuthGuardService],
        children: [
            {
                path: '',
                loadChildren: exports.projectsPath,
            }
        ]
    },
    {
        path: '',
        component: login_layout_component_1.LoginLayoutComponent,
        children: [
            {
                path: 'login',
                loadChildren: exports.loginPath
            }
        ]
    },
    { path: '**', redirectTo: '' }
];
let AppRoutingModule = class AppRoutingModule {
};
AppRoutingModule = __decorate([
    core_1.NgModule({
        imports: [
            router_1.RouterModule.forRoot(appRoutes, { useHash: false })
        ],
        exports: [
            router_1.RouterModule
        ]
    })
], AppRoutingModule);
exports.AppRoutingModule = AppRoutingModule;


/***/ }),

/***/ "./src/app/app.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const platform_browser_1 = __webpack_require__("./node_modules/@angular/platform-browser/esm2015/platform-browser.js");
const http_1 = __webpack_require__("./node_modules/@angular/common/esm2015/http.js");
const app_component_1 = __webpack_require__("./src/app/ontology.root/app.component.ts");
const shared_module_1 = __webpack_require__("./src/app/common/shared.module.ts");
const navbar_module_1 = __webpack_require__("./src/app/ontology.navbar/navbar.module.ts");
const footer_module_1 = __webpack_require__("./src/app/ontology.footer/footer.module.ts");
const cytoscape_service_1 = __webpack_require__("./src/app/services/cytoscape.service.ts");
const project_service_1 = __webpack_require__("./src/app/services/project.service.ts");
const conversation_service_1 = __webpack_require__("./src/app/services/conversation.service.ts");
const layout_service_1 = __webpack_require__("./src/app/ontology.layout/layout.service.ts");
const app_routing_module_1 = __webpack_require__("./src/app/app-routing.module.ts");
const project_new_module_1 = __webpack_require__("./src/app/ontology.project.new/project.new.module.ts");
const notification_module_1 = __webpack_require__("./src/app/ontology.notification/notification.module.ts");
const notification_service_1 = __webpack_require__("./src/app/services/notification.service.ts");
const http_handler_interceptor_1 = __webpack_require__("./src/app/http.handler.interceptor.ts");
const unfinished_topic_module_1 = __webpack_require__("./src/app/ontology.unfinished.topic/unfinished.topic.module.ts");
const auth_service_1 = __webpack_require__("./src/app/services/auth.service.ts");
const auth_guard_service_1 = __webpack_require__("./src/app/services/auth-guard.service.ts");
const home_layout_component_1 = __webpack_require__("./src/app/layouts/home-layout/home-layout.component.ts");
const login_layout_component_1 = __webpack_require__("./src/app/layouts/login-layout/login-layout.component.ts");
let AppModule = class AppModule {
};
AppModule = __decorate([
    core_1.NgModule({
        declarations: [
            app_component_1.AppComponent,
            home_layout_component_1.HomeLayoutComponent,
            login_layout_component_1.LoginLayoutComponent
        ],
        imports: [
            platform_browser_1.BrowserModule,
            http_1.HttpClientModule,
            shared_module_1.SharedModule,
            navbar_module_1.NavbarModule,
            footer_module_1.FooterModule,
            project_new_module_1.ProjectNewModule,
            notification_module_1.NotificationModule,
            unfinished_topic_module_1.UnfinishedTopicModule,
            platform_browser_1.BrowserModule,
            app_routing_module_1.AppRoutingModule //should be the latest module
        ],
        providers: [
            cytoscape_service_1.CytoscapeService,
            project_service_1.ProjectService,
            notification_service_1.NotificationService,
            conversation_service_1.ConversationService,
            auth_service_1.AuthService,
            auth_guard_service_1.AuthGuardService,
            layout_service_1.LayoutService,
            {
                provide: http_1.HTTP_INTERCEPTORS,
                useClass: http_handler_interceptor_1.HttpHandlerInterceptor,
                multi: true
            }
        ],
        bootstrap: [app_component_1.AppComponent]
    })
], AppModule);
exports.AppModule = AppModule;


/***/ }),

/***/ "./src/app/common/global.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
exports.BASE_API_URL = '/rest/api';
exports.cytoscape_style = [
    {
        "selector": "node",
        "style": {
            "font-family": "Source Sans Pro",
            "width": "32px",
            "height": "32px",
            "background-color": "#4FC3F7",
            "border-color": "#FFF",
            "border-width": 2,
            "content": "data(name)",
            "font-size": "13px"
        }
    },
    {
        "selector": "node#u",
        "style": {
            "background-color": "blue"
        }
    },
    {
        "selector": "edge",
        "style": {
            "width": 2,
            "line-color": "#B3E5FC",
            "target-arrow-color": "#B3E5FC",
            "source-arrow-color": "#B3E5FC",
            "curve-style": "bezier",
            "target-arrow-shape": "vee",
            "font-size": "12px",
            "content": "data(name)"
        }
    },
    {
        "selector": "edge#pu",
        "style": {
            "source-arrow-color": "#B6B6B6",
            "target-arrow-shape": "vee",
            "mid-target-arrow-shape": "vee",
            "source-arrow-shape": "vee"
        }
    },
    {
        "selector": ".eh-handle",
        "style": {
            "background-color": "red",
            "width": 12,
            "height": 12,
            "shape": "ellipse",
            "overlay-opacity": 0,
            "border-width": 1,
            "border-opacity": 0
        }
    },
    {
        "selector": ".eh-hover",
        "style": {
            "background-color": "red"
        }
    },
    {
        "selector": ".eh-source",
        "style": {
            "border-width": 2,
            "border-color": "red"
        }
    },
    {
        "selector": ".eh-target",
        "style": {
            "border-width": 2,
            "border-color": "red"
        }
    },
    {
        "selector": ".eh-preview, .eh-ghost-edge",
        "style": {
            "background-color": "red",
            "line-color": "red",
            "target-arrow-color": "red"
        }
    },
    {
        "selector": ".node-on-hover-highlight",
        "style": {
            "background-color": "#4FC3F7",
            "border-color": "#E1F5FE",
            "border-width": "2px"
        }
    },
    {
        "selector": ".node-on-hover-edge-highlight",
        "style": {
            "line-color": "#0081CB",
            "target-arrow-color": "#0081CB",
            "mid-target-arrow-color": "#0081CB",
            "source-arrow-color": "#0081CB"
        }
    },
    {
        "selector": ".node-on-select-highlight",
        "style": {
            "height": "32px",
            "width": "32px",
            "background-color": "#0081CB",
            "border-color": "#FFFFFF",
            "border-width": "2px"
        }
    },
    {
        "selector": ".node-on-select-edge-highlight",
        "style": {
            "line-color": "#4FC3F7",
            "target-arrow-color": "#4FC3F7",
            "mid-target-arrow-color": "#4FC3F7",
            "source-arrow-color": "#4FC3F7"
        }
    },
    {
        "selector": ".node-bigger",
        "style": {
            "background-color": "blue",
            "width": 30,
            "height": 30
        }
    },
    {
        "selector": ".attribute",
        "style": {
            "font-family": "Source Sans Pro",
            "font-size": "12px",
            "background-color": "#FF8456",
            "height": "16px",
            "width": "16px",
            "border-width": 0
        }
    },
    {
        "selector": ":parent",
        "style": {
            "background-color": "#E8F4F9",
            "background-opacity": 0.7
        }
    },
    {
        "selector": "node.cy-expand-collapse-collapsed-node",
        "style": {
            "background-fit": "cover",
        }
    }
];


/***/ }),

/***/ "./src/app/common/global_local.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
exports.BASE_API_URL = '/rest/api';
exports.cytoscape_style = [
    {
        "selector": "node",
        "style": {
            "font-family": "Source Sans Pro",
            "width": "32px",
            "height": "32px",
            "background-color": "#4FC3F7",
            "border-color": "#FFF",
            "border-width": 2,
            "content": "data(name)",
            "font-size": "13px"
        }
    },
    {
        "selector": "node#u",
        "style": {
            "background-color": "blue"
        }
    },
    {
        "selector": "edge",
        "style": {
            "width": 2,
            "line-color": "#B3E5FC",
            "target-arrow-color": "#B3E5FC",
            "source-arrow-color": "#B3E5FC",
            "curve-style": "bezier",
            "target-arrow-shape": "vee",
            "font-size": "12px",
            "content": "data(name)"
        }
    },
    {
        "selector": "edge#pu",
        "style": {
            "source-arrow-color": "#B6B6B6",
            "target-arrow-shape": "vee",
            "mid-target-arrow-shape": "vee",
            "source-arrow-shape": "vee"
        }
    },
    {
        "selector": ".eh-handle",
        "style": {
            "background-color": "red",
            "width": 12,
            "height": 12,
            "shape": "ellipse",
            "overlay-opacity": 0,
            "border-width": 1,
            "border-opacity": 0
        }
    },
    {
        "selector": ".eh-hover",
        "style": {
            "background-color": "red"
        }
    },
    {
        "selector": ".eh-source",
        "style": {
            "border-width": 2,
            "border-color": "red"
        }
    },
    {
        "selector": ".eh-target",
        "style": {
            "border-width": 2,
            "border-color": "red"
        }
    },
    {
        "selector": ".eh-preview, .eh-ghost-edge",
        "style": {
            "background-color": "red",
            "line-color": "red",
            "target-arrow-color": "red"
        }
    },
    {
        "selector": ".node-on-hover-highlight",
        "style": {
            "background-color": "#4FC3F7",
            "border-color": "#E1F5FE",
            "border-width": "2px"
        }
    },
    {
        "selector": ".node-on-hover-edge-highlight",
        "style": {
            "line-color": "#0081CB",
            "target-arrow-color": "#0081CB",
            "mid-target-arrow-color": "#0081CB",
            "source-arrow-color": "#0081CB"
        }
    },
    {
        "selector": ".node-on-select-highlight",
        "style": {
            "height": "32px",
            "width": "32px",
            "background-color": "#0081CB",
            "border-color": "#FFFFFF",
            "border-width": "2px"
        }
    },
    {
        "selector": ".node-on-select-edge-highlight",
        "style": {
            "line-color": "#4FC3F7",
            "target-arrow-color": "#4FC3F7",
            "mid-target-arrow-color": "#4FC3F7",
            "source-arrow-color": "#4FC3F7"
        }
    },
    {
        "selector": ".node-bigger",
        "style": {
            "background-color": "blue",
            "width": 30,
            "height": 30
        }
    },
    {
        "selector": ".attribute",
        "style": {
            "font-family": "Source Sans Pro",
            "font-size": "12px",
            "background-color": "#FF8456",
            "height": "16px",
            "width": "16px",
            "border-width": 0
        }
    },
    {
        "selector": ":parent",
        "style": {
            "background-color": "#E8F4F9",
            "background-opacity": 0.7
        }
    },
    {
        "selector": "node.cy-expand-collapse-collapsed-node",
        "style": {
            "background-fit": "cover",
        }
    }
];


/***/ }),

/***/ "./src/app/common/global_remote.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
exports.BASE_API_URL = '/rest/api';
exports.cytoscape_style = [
    {
        "selector": "node",
        "style": {
            "font-family": "Source Sans Pro",
            "width": "32px",
            "height": "32px",
            "background-color": "#E778EE",
            "border-color": "#FFF",
            "border-width": 2,
            "content": "data(name)",
            "font-size": "13px"
        }
    },
    {
        "selector": "node#u",
        "style": {
            "background-color": "blue"
        }
    },
    {
        "selector": "edge",
        "style": {
            "width": 2,
            "line-color": "#B3E5FC",
            "target-arrow-color": "#B3E5FC",
            "source-arrow-color": "#B3E5FC",
            "curve-style": "bezier",
            "target-arrow-shape": "vee",
            "font-size": "12px",
            "content": "data(name)"
        }
    },
    {
        "selector": "edge#pu",
        "style": {
            "source-arrow-color": "#B6B6B6",
            "target-arrow-shape": "vee",
            "mid-target-arrow-shape": "vee",
            "source-arrow-shape": "vee"
        }
    },
    {
        "selector": ".eh-handle",
        "style": {
            "background-color": "red",
            "width": 12,
            "height": 12,
            "shape": "ellipse",
            "overlay-opacity": 0,
            "border-width": 1,
            "border-opacity": 0
        }
    },
    {
        "selector": ".eh-hover",
        "style": {
            "background-color": "red"
        }
    },
    {
        "selector": ".eh-source",
        "style": {
            "border-width": 2,
            "border-color": "red"
        }
    },
    {
        "selector": ".eh-target",
        "style": {
            "border-width": 2,
            "border-color": "red"
        }
    },
    {
        "selector": ".eh-preview, .eh-ghost-edge",
        "style": {
            "background-color": "red",
            "line-color": "red",
            "target-arrow-color": "red"
        }
    },
    {
        "selector": ".node-on-hover-highlight",
        "style": {
            "background-color": "#4FC3F7",
            "border-color": "#E1F5FE",
            "border-width": "2px"
        }
    },
    {
        "selector": ".node-on-hover-edge-highlight",
        "style": {
            "line-color": "#0081CB",
            "target-arrow-color": "#0081CB",
            "mid-target-arrow-color": "#0081CB",
            "source-arrow-color": "#0081CB"
        }
    },
    {
        "selector": ".node-on-select-highlight",
        "style": {
            "height": "32px",
            "width": "32px",
            "background-color": "#0081CB",
            "border-color": "#FFFFFF",
            "border-width": "2px"
        }
    },
    {
        "selector": ".node-on-select-edge-highlight",
        "style": {
            "line-color": "#4FC3F7",
            "target-arrow-color": "#4FC3F7",
            "mid-target-arrow-color": "#4FC3F7",
            "source-arrow-color": "#4FC3F7"
        }
    },
    {
        "selector": ".node-bigger",
        "style": {
            "background-color": "blue",
            "width": 30,
            "height": 30
        }
    },
    {
        "selector": ".attribute",
        "style": {
            "font-family": "Source Sans Pro",
            "font-size": "12px",
            "background-color": "#FF8456",
            "height": "16px",
            "width": "16px",
            "border-width": 0
        }
    },
    {
        "selector": ":parent",
        "style": {
            "background-color": "#E8F4F9",
            "background-opacity": 0.7
        }
    },
    {
        "selector": "node.cy-expand-collapse-collapsed-node",
        "style": {
            "background-fit": "cover",
        }
    }
];


/***/ }),

/***/ "./src/app/common/shared.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const common_1 = __webpack_require__("./node_modules/@angular/common/esm2015/common.js");
const forms_1 = __webpack_require__("./node_modules/@angular/forms/esm2015/forms.js");
let SharedModule = class SharedModule {
};
SharedModule = __decorate([
    core_1.NgModule({
        imports: [common_1.CommonModule],
        declarations: [],
        exports: [common_1.CommonModule, forms_1.FormsModule, forms_1.ReactiveFormsModule]
    })
], SharedModule);
exports.SharedModule = SharedModule;


/***/ }),

/***/ "./src/app/http.handler.interceptor.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const Observable_1 = __webpack_require__("./node_modules/rxjs/_esm2015/Observable.js");
const notification_service_1 = __webpack_require__("./src/app/services/notification.service.ts");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
__webpack_require__("./node_modules/rxjs/_esm2015/add/operator/catch.js");
__webpack_require__("./node_modules/rxjs/_esm2015/add/observable/throw.js");
const auth_service_1 = __webpack_require__("./src/app/services/auth.service.ts");
let HttpHandlerInterceptor = class HttpHandlerInterceptor {
    constructor(notificationService, authService, router) {
        this.notificationService = notificationService;
        this.authService = authService;
        this.router = router;
    }
    intercept(req, next) {
        let reqHeader = req;
        if (!req.url.includes("/users/login")) {
            reqHeader = req.clone({
                headers: req.headers.set('Authorization', sessionStorage.getItem('otyToken'))
            });
        }
        return next.handle(reqHeader)
            .catch((error) => {
            if (this.router.url !== '/login' && (error.status === 401 || error.status === 504)) {
                this.authService.deleteSession();
                this.router.navigate(['/login']);
            }
            return Observable_1.Observable.throw(error);
        })
            .map((res) => {
            this.errorHandler(res.body);
            return res;
        });
    }
    errorHandler(body) {
        if (body && body.error) {
            this.showErrorMessage('Error: ' + body.error.message);
        }
    }
    showErrorMessage(message) {
        this.notificationService.showNotification('error', message);
    }
};
HttpHandlerInterceptor = __decorate([
    core_1.Injectable(),
    __metadata("design:paramtypes", [notification_service_1.NotificationService,
        auth_service_1.AuthService,
        router_1.Router])
], HttpHandlerInterceptor);
exports.HttpHandlerInterceptor = HttpHandlerInterceptor;


/***/ }),

/***/ "./src/app/layouts/home-layout/home-layout.component.main.scss":
/***/ (function(module, exports) {

module.exports = "#oty-main-content {\n  height: calc(100% - 3px); }\n\noty-notification {\n  height: auto; }\n"

/***/ }),

/***/ "./src/app/layouts/home-layout/home-layout.component.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
let HomeLayoutComponent = class HomeLayoutComponent {
};
HomeLayoutComponent = __decorate([
    core_1.Component({
        selector: 'oty-home-layout',
        template: `
    <oty-navbar></oty-navbar>
    <oty-notification></oty-notification>
    <div id="oty-main-content">
      <router-outlet></router-outlet>
    </div>
  `,
        styles: [__webpack_require__("./src/app/layouts/home-layout/home-layout.component.main.scss")]
    })
], HomeLayoutComponent);
exports.HomeLayoutComponent = HomeLayoutComponent;
/* To-do: Figure out why CSS styles are not being applied to the layout */ 


/***/ }),

/***/ "./src/app/layouts/login-layout/login-layout.component.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
let LoginLayoutComponent = class LoginLayoutComponent {
};
LoginLayoutComponent = __decorate([
    core_1.Component({
        selector: 'oty-login-layout',
        template: `    
    <oty-notification></oty-notification>
    <div id="oty-main-content">
      <router-outlet></router-outlet>
    </div>
  `,
        styles: []
    })
], LoginLayoutComponent);
exports.LoginLayoutComponent = LoginLayoutComponent;


/***/ }),

/***/ "./src/app/ontology.footer/footer.component.html":
/***/ (function(module, exports) {

module.exports = "<footer>\n  this is a footer.\n</footer>\n"

/***/ }),

/***/ "./src/app/ontology.footer/footer.component.main.scss":
/***/ (function(module, exports) {

module.exports = "footer {\n  height: 50px;\n  line-height: 50px;\n  background-color: #062e79;\n  text-align: center;\n  font-size: 16px;\n  color: #FFF; }\n"

/***/ }),

/***/ "./src/app/ontology.footer/footer.component.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
let FooterComponent = class FooterComponent {
};
FooterComponent = __decorate([
    core_1.Component({
        selector: 'oty-footer',
        template: __webpack_require__("./src/app/ontology.footer/footer.component.html"),
        styles: [__webpack_require__("./src/app/ontology.footer/footer.component.main.scss")]
    })
], FooterComponent);
exports.FooterComponent = FooterComponent;


/***/ }),

/***/ "./src/app/ontology.footer/footer.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const footer_component_1 = __webpack_require__("./src/app/ontology.footer/footer.component.ts");
let FooterModule = class FooterModule {
};
FooterModule = __decorate([
    core_1.NgModule({
        imports: [],
        exports: [footer_component_1.FooterComponent],
        declarations: [footer_component_1.FooterComponent]
    })
], FooterModule);
exports.FooterModule = FooterModule;


/***/ }),

/***/ "./src/app/ontology.layout/layout.service.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
let LayoutService = class LayoutService {
    constructor() {
        this.toggleDetailPanelEvent = new core_1.EventEmitter();
        this.onSetDetailData = new core_1.EventEmitter();
    }
    toggleDetailPanel(isShow, isRightSide) {
        let event = {
            'isShow': isShow,
            'isRightSide': isRightSide
        };
        this.toggleDetailPanelEvent.emit(event);
    }
    setDetailData(graph) {
        this.onSetDetailData.emit(graph);
    }
};
LayoutService = __decorate([
    core_1.Injectable()
], LayoutService);
exports.LayoutService = LayoutService;


/***/ }),

/***/ "./src/app/ontology.navbar/navbar.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"nav\">\n  <div id=\"navbar\" ui-view=\"\" class=\"tc-global-header-container\" >\n    <div class=\"tc-global-header\">\n      <a class=\"oty-icon-header\" routerLink='/home'>\n        <i class=\"oty-icon-logo\"></i>\n      </a>\n      <div class=\"tc-icon-tray\" [ngClass]=\"{'active': selectedIndex!==-1}\">\n        <ul class=\"tc-global-header-nav\">\n          <li class=\"tc-header-nav-item oty-icon-signout\" (click)=\"logout()\"></li>\n          <li class=\"tc-header-nav-item profile\">\n            <span class=\"tc-profile-name\">{{username}}</span>\n          </li>\n        </ul>\n      </div>\n    </div>\n    <div class=\"tc-container-background\" [ngStyle]=\"{display: selectedIndex!==-1?'block':'none'}\" (click)=\"resetIndex()\"></div>\n    <div class=\"tc-panel-background\" [ngStyle]=\"{display: selectedIndex!==-1?'block':'none'}\"></div>\n    <div class=\"tc-panel-container\" [ngClass]=\"{'slide-in': selectedIndex!==-1}\">\n      <div class=\"tc-panel-body\"></div>\n      <div class=\"tc-panel-footer\"></div>\n    </div>\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.navbar/navbar.component.main.scss":
/***/ (function(module, exports) {

module.exports = ".oty-icon-header {\n  padding-left: 25px;\n  color: transparent; }\n  .oty-icon-header .oty-icon-logo {\n    font-size: 40px;\n    color: #FFFFFF; }\n  .oty-icon-header span {\n    position: fixed; }\n  .tc-icon-tray {\n  width: auto; }\n  .tc-container-background {\n  z-index: 9998 !important; }\n  .tc-global-header-nav {\n  margin: 0;\n  cursor: inherit !important; }\n  .tc-global-header-nav li.tc-header-nav-item.oty-icon-signout {\n    margin-left: 10px;\n    width: 40px !important;\n    cursor: pointer; }\n"

/***/ }),

/***/ "./src/app/ontology.navbar/navbar.component.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const auth_service_1 = __webpack_require__("./src/app/services/auth.service.ts");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
let NavbarComponent = class NavbarComponent {
    constructor(authService, router) {
        this.authService = authService;
        this.router = router;
        this.selectedIndex = -1;
    }
    ngOnInit() {
        this.authService
            .usernameEvent
            .subscribe((username) => {
            this.username = username;
        });
        this.username = this.authService.getUsername();
    }
    setSelectedIndex(index) {
        this.selectedIndex = index;
    }
    resetIndex() {
        if (this.selectedIndex != -1)
            this.selectedIndex = -1;
    }
    logout() {
        this.authService.logout().subscribe((isLogout) => {
            if (isLogout) {
                this.resetIndex();
                this.router.navigate(['/login']);
            }
        });
    }
};
NavbarComponent = __decorate([
    core_1.Component({
        selector: 'oty-navbar',
        template: __webpack_require__("./src/app/ontology.navbar/navbar.component.html"),
        styles: [__webpack_require__("./src/app/ontology.navbar/navbar.component.main.scss")]
    }),
    __metadata("design:paramtypes", [auth_service_1.AuthService, router_1.Router])
], NavbarComponent);
exports.NavbarComponent = NavbarComponent;


/***/ }),

/***/ "./src/app/ontology.navbar/navbar.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const common_1 = __webpack_require__("./node_modules/@angular/common/esm2015/common.js");
const navbar_component_1 = __webpack_require__("./src/app/ontology.navbar/navbar.component.ts");
let NavbarModule = class NavbarModule {
};
NavbarModule = __decorate([
    core_1.NgModule({
        imports: [common_1.CommonModule, router_1.RouterModule],
        exports: [navbar_component_1.NavbarComponent],
        declarations: [navbar_component_1.NavbarComponent]
    })
], NavbarModule);
exports.NavbarModule = NavbarModule;


/***/ }),

/***/ "./src/app/ontology.notification/notification.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"oty-wrapper-notification\" *ngIf=\"isVisible\">\n  <div class=\"tc-notifications\" [ngClass]=\"'tc-notifications-'+type\">\n    <span [ngClass]=\"'tc-'+type+'-close'\" class=\"tc-info-close}}\" (click)=\"close()\"></span>\n    <div class=\"tc-notifications-message\">{{message}}</div>\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.notification/notification.component.scss":
/***/ (function(module, exports) {

module.exports = ".oty-wrapper-notification {\n  position: fixed;\n  margin-top: 4%;\n  margin-left: 30%;\n  z-index: 9999 !important;\n  -webkit-animation: fadeOut 4s linear forwards;\n          animation: fadeOut 4s linear forwards; }\n  .oty-wrapper-notification .tc-notifications {\n    height: auto;\n    width: 94%; }\n  .oty-wrapper-notification .tc-notifications .tc-notifications-message {\n      width: calc(100% - 30px); }\n  @-webkit-keyframes fadeOut {\n  0% {\n    opacity: 0; }\n  10% {\n    opacity: 1; }\n  90% {\n    opacity: 1;\n    -webkit-transform: translateY(0px);\n            transform: translateY(0px); }\n  99% {\n    opacity: 0;\n    -webkit-transform: translateY(-30px);\n            transform: translateY(-30px); }\n  100% {\n    opacity: 0; } }\n  @keyframes fadeOut {\n  0% {\n    opacity: 0; }\n  10% {\n    opacity: 1; }\n  90% {\n    opacity: 1;\n    -webkit-transform: translateY(0px);\n            transform: translateY(0px); }\n  99% {\n    opacity: 0;\n    -webkit-transform: translateY(-30px);\n            transform: translateY(-30px); }\n  100% {\n    opacity: 0; } }\n"

/***/ }),

/***/ "./src/app/ontology.notification/notification.component.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const notification_service_1 = __webpack_require__("./src/app/services/notification.service.ts");
var notificationType;
(function (notificationType) {
    notificationType["success"] = "success";
    notificationType["error"] = "error";
    notificationType["info"] = "info";
})(notificationType = exports.notificationType || (exports.notificationType = {}));
let NotificationComponent = class NotificationComponent {
    constructor(notificationService) {
        this.notificationService = notificationService;
        this.type = notificationType.info;
        this.message = 'Default';
        this.isVisible = false;
    }
    ngOnInit() {
        this.notificationService
            .notificationEvent
            .subscribe((res) => {
            this.showNotification(res.type, res.message);
        });
    }
    close() {
        this.isVisible = false;
    }
    showNotification(type, message) {
        if (!this.isVisible) {
            switch (type) {
                case notificationType.success:
                    this.type = notificationType.success;
                    break;
                case notificationType.info:
                    this.type = notificationType.info;
                    break;
                case notificationType.error:
                    this.type = notificationType.error;
                    break;
            }
            this.message = message;
            this.isVisible = true;
        }
    }
};
NotificationComponent = __decorate([
    core_1.Component({
        selector: 'oty-notification',
        template: __webpack_require__("./src/app/ontology.notification/notification.component.html"),
        styles: [__webpack_require__("./src/app/ontology.notification/notification.component.scss")]
    }),
    __metadata("design:paramtypes", [notification_service_1.NotificationService])
], NotificationComponent);
exports.NotificationComponent = NotificationComponent;


/***/ }),

/***/ "./src/app/ontology.notification/notification.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const common_1 = __webpack_require__("./node_modules/@angular/common/esm2015/common.js");
const notification_component_1 = __webpack_require__("./src/app/ontology.notification/notification.component.ts");
let NotificationModule = class NotificationModule {
};
NotificationModule = __decorate([
    core_1.NgModule({
        imports: [
            common_1.CommonModule
        ],
        exports: [
            notification_component_1.NotificationComponent
        ],
        declarations: [
            notification_component_1.NotificationComponent
        ]
    })
], NotificationModule);
exports.NotificationModule = NotificationModule;


/***/ }),

/***/ "./src/app/ontology.project.new/project.new-routing.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const project_new_component_1 = __webpack_require__("./src/app/ontology.project.new/project.new.component.ts");
const auth_guard_service_1 = __webpack_require__("./src/app/services/auth-guard.service.ts");
const routes = [
    { path: 'project-new', component: project_new_component_1.ProjectNewComponent, canActivate: [auth_guard_service_1.AuthGuardService] }
];
let ProjectNewRoutingModule = class ProjectNewRoutingModule {
};
ProjectNewRoutingModule = __decorate([
    core_1.NgModule({
        imports: [router_1.RouterModule.forChild(routes)],
        exports: [router_1.RouterModule]
    })
], ProjectNewRoutingModule);
exports.ProjectNewRoutingModule = ProjectNewRoutingModule;


/***/ }),

/***/ "./src/app/ontology.project.new/project.new.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"gdb-project-new-background\">\n  <div class=\"conversation-bot-wrapper\">\n    <div id=\"oty-bot-icon\"><img src=\"../../assets/images/bot-icon.png\" alt=\"profile-bot\"></div>\n    <div class=\"message-section\">\n      <div class=\"rectangle-start\">Hi again!</div>\n      <div class=\"rectangle-final\"><span>Let's get started! What would you like to name your new project?</span></div>\n    </div>\n    <form [formGroup]=\"elemForm\" (ngSubmit)=\"onSubmit(elemForm)\">\n\n      <input type=\"text\" pInputText placeholder=\"Type project name here\" formControlName=\"name\" required/>\n      <button type=\"submit\" [disabled]=\"!elemForm.valid\">\n        <img src=\"../../assets/images/send.png\" alt=\"send\">\n      </button>\n    </form>\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.project.new/project.new.component.scss":
/***/ (function(module, exports) {

module.exports = ".gdb-project-new-background {\n  position: fixed;\n  width: 100%;\n  height: 100%;\n  font-size: 18px;\n  letter-spacing: 0.3px;\n  line-height: 19px;\n  left: 0;\n  top: 0;\n  background: linear-gradient(180.21deg, #0081cb 0%, #FFFFFF 100%); }\n  .gdb-project-new-background .conversation-bot-wrapper {\n    width: 40%;\n    height: 50%;\n    position: fixed;\n    top: 25%;\n    left: 40%; }\n  .gdb-project-new-background .conversation-bot-wrapper #oty-bot-icon {\n      float: left;\n      display: inline-block;\n      height: 104px;\n      width: 104px;\n      border-radius: 50% 50%; }\n  .gdb-project-new-background .conversation-bot-wrapper #oty-bot-icon img {\n        width: 100%;\n        height: 100%; }\n  .gdb-project-new-background .conversation-bot-wrapper .message-section {\n      display: inline-block; }\n  .gdb-project-new-background .conversation-bot-wrapper .message-section .rectangle-start,\n      .gdb-project-new-background .conversation-bot-wrapper .message-section .rectangle-final {\n        border: 1px solid #DEDEDE;\n        opacity: 0.9;\n        background-color: #FFFFFF;\n        padding: 2% 2% 2% 3%; }\n  .gdb-project-new-background .conversation-bot-wrapper .message-section .rectangle-start {\n        margin: 1% 1% 8px 11px;\n        height: 38px;\n        width: 107px;\n        border-radius: 0 15px 15px 15px; }\n  .gdb-project-new-background .conversation-bot-wrapper .message-section .rectangle-final {\n        margin: 3% 1% 37px 11px;\n        height: 59px;\n        width: 324px;\n        border-radius: 15px; }\n  .gdb-project-new-background .conversation-bot-wrapper form {\n      width: 498px;\n      position: absolute;\n      left: 0; }\n  .gdb-project-new-background .conversation-bot-wrapper form input {\n        color: #727272;\n        font-size: 16px;\n        width: 441px;\n        height: 66px;\n        border-radius: 15px;\n        border: 0;\n        padding-left: 5%; }\n  .gdb-project-new-background .conversation-bot-wrapper form input:focus {\n          outline: none; }\n  .gdb-project-new-background .conversation-bot-wrapper form button {\n        color: white;\n        position: absolute;\n        background-color: #2694D3;\n        border-radius: 50%;\n        border: 0;\n        top: 9px;\n        right: 0;\n        height: 48px;\n        width: 48px;\n        padding-left: 13px; }\n  .gdb-project-new-background .conversation-bot-wrapper form button:focus {\n          outline: none; }\n  .gdb-project-new-background .conversation-bot-wrapper form button:hover {\n          background-color: #0081cb; }\n  .gdb-project-new-background .conversation-bot-wrapper form button:active {\n          background-color: #E0F0F9; }\n  .gdb-project-new-background .conversation-bot-wrapper form button img {\n          height: 19.56px;\n          width: 23px; }\n"

/***/ }),

/***/ "./src/app/ontology.project.new/project.new.component.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const project_service_1 = __webpack_require__("./src/app/services/project.service.ts");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const notification_service_1 = __webpack_require__("./src/app/services/notification.service.ts");
const forms_1 = __webpack_require__("./node_modules/@angular/forms/esm2015/forms.js");
let ProjectNewComponent = class ProjectNewComponent {
    constructor(projectService, notificationService, router, fb) {
        this.projectService = projectService;
        this.notificationService = notificationService;
        this.router = router;
        this.formBuilder = fb;
    }
    ngOnInit() {
        this.elemForm = this.formBuilder.group({
            'name': ["", forms_1.Validators.required],
        });
    }
    onSubmit(f) {
        const projectName = f.value.name;
        this.projectService.createProject(projectName)
            .subscribe((res) => {
            if (res.error) {
                this.notificationService.showNotification('error', res.error.message);
            }
            else {
                this.notificationService.showNotification('success', 'Project created');
                this.router.navigate(['/home']);
            }
        });
    }
};
ProjectNewComponent = __decorate([
    core_1.Component({
        selector: 'gdb-project-new',
        template: __webpack_require__("./src/app/ontology.project.new/project.new.component.html"),
        styles: [__webpack_require__("./src/app/ontology.project.new/project.new.component.scss")]
    }),
    __param(3, core_1.Inject(forms_1.FormBuilder)),
    __metadata("design:paramtypes", [project_service_1.ProjectService,
        notification_service_1.NotificationService,
        router_1.Router,
        forms_1.FormBuilder])
], ProjectNewComponent);
exports.ProjectNewComponent = ProjectNewComponent;


/***/ }),

/***/ "./src/app/ontology.project.new/project.new.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const common_1 = __webpack_require__("./node_modules/@angular/common/esm2015/common.js");
const project_new_component_1 = __webpack_require__("./src/app/ontology.project.new/project.new.component.ts");
const project_new_routing_module_1 = __webpack_require__("./src/app/ontology.project.new/project.new-routing.module.ts");
const shared_module_1 = __webpack_require__("./src/app/common/shared.module.ts");
let ProjectNewModule = class ProjectNewModule {
};
ProjectNewModule = __decorate([
    core_1.NgModule({
        imports: [
            common_1.CommonModule,
            project_new_routing_module_1.ProjectNewRoutingModule,
            shared_module_1.SharedModule
        ],
        declarations: [project_new_component_1.ProjectNewComponent]
    })
], ProjectNewModule);
exports.ProjectNewModule = ProjectNewModule;


/***/ }),

/***/ "./src/app/ontology.root/app.component.html":
/***/ (function(module, exports) {

module.exports = "<router-outlet></router-outlet>\n"

/***/ }),

/***/ "./src/app/ontology.root/app.component.scss":
/***/ (function(module, exports) {

module.exports = "#oty-main-content {\n  height: 100%;\n  margin: 0 auto; }\n"

/***/ }),

/***/ "./src/app/ontology.root/app.component.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
let AppComponent = class AppComponent {
};
AppComponent = __decorate([
    core_1.Component({
        selector: 'oty-root',
        template: __webpack_require__("./src/app/ontology.root/app.component.html"),
        styles: [__webpack_require__("./src/app/ontology.root/app.component.scss")]
    })
], AppComponent);
exports.AppComponent = AppComponent;


/***/ }),

/***/ "./src/app/ontology.unfinished.topic/unfinished.topic.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"oty-unfinished-topic-container\" [ngClass]=\"{'expand': isShowSidebar}\">\n  <div class=\"oty-unfinished-topic-content\">\n    <div class=\"oty-search-section\">\n      <p><span *ngIf=\"!showSearchInput\">TODO LIST</span>\n        <input type=\"text\" (input)=\"onSearchChange($event.target.value)\"\n               title=\"search topic\"\n               type=\"text\"\n               *ngIf=\"showSearchInput\"\n               placeholder=\"Search\"\n               #search >\n        <button class=\"oty-search-button\" (click)=\"focusInput()\">\n          <i class=\"oty-icon-search\"></i>\n        </button>\n      </p>\n    </div>\n    <div>\n      <div class=\"oty-bubble-speech\" *ngFor=\"let topic of unfinishedTopics; let i = index\"\n           (onDragStart)=\"onDragStart($event, topic, i)\" (onDragEnd)=\"onDragEnd($event, topic, i)\"  pDroppable=\"oty-unfinished-topic\"\n           pDraggable=\"oty-unfinished-topic\">\n        <span title=\"{{topic.lastModified | date:'short'}}\">{{topic.title}}</span>\n      </div>\n    </div>\n\n    <div class=\"area-unselected conv-area-unselected\" [ngClass]=\"{'conv-area-selected':isOverCurrentConvArea}\"\n         pDroppable=\"oty-unfinished-topic\" *ngIf=\"isShowDropArea\"\n         (onDragEnter)=\"isOverCurrentConvArea = true\" (onDragLeave)=\"isOverCurrentConvArea = false\"\n         (onDrop)=\"setCurrentConversation($event)\"\n    >continue topic</div>\n    <div class=\"area-unselected delete-area-unselected\" [ngClass]=\"{'delete-area-selected':isOverDeleteConvArea}\"\n         pDroppable=\"oty-unfinished-topic\" *ngIf=\"isShowDropArea\"\n         (onDragEnter)=\"isOverDeleteConvArea = true\" (onDragLeave)=\"isOverDeleteConvArea = false\"\n         (onDrop)=\"deleteConversation($event)\">\n      <i class=\"oty-icon-trashcan\"></i>Delete\n    </div>\n  </div>\n  <div class=\"oty-selector\" (click)=\"openSideBar()\">\n    <img src=\"../../assets/images/open-tab.png\">\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.unfinished.topic/unfinished.topic.component.scss":
/***/ (function(module, exports) {

module.exports = ".oty-unfinished-topic-container {\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  height: 100%;\n  position: absolute;\n  text-align: center;\n  z-index: 2; }\n  .oty-unfinished-topic-container.expand {\n    width: 100%; }\n  .oty-unfinished-topic-container.expand .oty-unfinished-topic-content {\n      display: block; }\n  .oty-unfinished-topic-container.expand .oty-selector {\n      right: 0; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content {\n    background-color: #F4F4F4;\n    color: #727272;\n    height: 100%;\n    overflow: hidden;\n    position: relative;\n    width: 100%;\n    display: none; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-search-section input {\n      border: 0;\n      border-radius: 15px;\n      width: calc(100% - 100px);\n      height: 26px;\n      padding-left: 2%; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-search-section input:focus {\n        outline: none; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-search-section > p {\n      font-size: 12px;\n      font-weight: bold;\n      letter-spacing: 1.08px;\n      line-height: 15px;\n      margin-bottom: 24px;\n      margin-top: 24px; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-search-section > p span {\n        display: inline-block;\n        margin-top: 4px; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-search-section .oty-search-button {\n      background-color: transparent;\n      border: 0;\n      display: inline-block;\n      position: absolute;\n      right: 24px;\n      top: 28px;\n      font-size: 13px; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-search-section .oty-search-button:focus {\n        outline: none; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-search-section .oty-search-button:active {\n        color: #B6B6B6; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-search-section .oty-search-button .oty-icon-search {\n        font-weight: 600; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content > div:nth-child(2) {\n      overflow: auto;\n      height: calc(100% - 48px - 150px); }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-bubble-speech {\n      background: #FFFFFF;\n      border-radius: 3px;\n      margin-bottom: 8px;\n      margin-left: 16px;\n      margin-right: 16px;\n      min-height: 48px;\n      text-align: left;\n      padding: 10px;\n      cursor: pointer; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .area-unselected {\n      border: 2px dashed;\n      font-size: 16px;\n      height: 48px;\n      padding: 10px; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .conv-area-unselected {\n      border-color: #0081cb;\n      color: #0081cb;\n      position: absolute;\n      bottom: 60px;\n      width: 100%; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .delete-area-unselected {\n      border-color: #D0021B;\n      bottom: 0;\n      color: #D0021B;\n      position: absolute;\n      width: 100%; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .delete-area-selected {\n      background: #FFB6BF; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .conv-area-selected {\n      background: #E0F0F9; }\n  .oty-unfinished-topic-container .oty-unfinished-topic-content .oty-icon-trashcan {\n      font-size: 24px;\n      left: 35%;\n      position: absolute; }\n  .oty-unfinished-topic-container .oty-selector {\n    position: absolute;\n    cursor: pointer;\n    background: none;\n    width: 20px;\n    height: 35px; }\n  .oty-unfinished-topic-container .oty-selector img {\n      width: 15px; }\n"

/***/ }),

/***/ "./src/app/ontology.unfinished.topic/unfinished.topic.component.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const conversation_service_1 = __webpack_require__("./src/app/services/conversation.service.ts");
let UnifinishedTopicComponent = class UnifinishedTopicComponent {
    constructor(conversationService) {
        this.conversationService = conversationService;
        this.isShowSidebar = false;
        this.showSearchInput = false;
        this.isOverCurrentConvArea = false;
        this.isOverDeleteConvArea = false;
        this.isShowDropArea = false;
    }
    onDragStart(event, rerunTopic, i) {
        this.showOptions();
        event.dataTransfer.setData('resumeTopic', JSON.stringify(rerunTopic));
        event.dataTransfer.setData('indexCurrentTopic', i);
    }
    onDragEnd(event, rerunTopic, i) {
        this.resetColorAreas();
        this.hideOptions();
    }
    openSideBar() {
        this.isShowSidebar = !this.isShowSidebar;
        this.unfinishedTopics = this.conversationService.getUnfinishedTopics();
    }
    showOptions() {
        this.isShowDropArea = true;
    }
    hideOptions() {
        this.isShowDropArea = false;
    }
    getCurrentConversation(event) {
        return JSON.parse(event.dataTransfer.getData('resumeTopic'));
    }
    blurInput() {
        this.unfinishedTopics = this.conversationService.getUnfinishedTopics();
        this.showSearchInput = false;
    }
    focusInput() {
        setTimeout(() => {
            this.inputEl.nativeElement.focus();
        }, 0);
        this.showSearchInput = true;
    }
    onSearchChange(value) {
        this.unfinishedTopics = this.conversationService.getUnfinishedTopics().filter((item) => {
            return item.title.toLowerCase().includes(value.toLowerCase());
        });
    }
    deleteConversation(event) {
        const index = event.dataTransfer.getData('indexCurrentTopic');
        this.unfinishedTopics.splice(index, 1);
        this.resetColorAreas();
        this.blurInput();
    }
    setCurrentConversation(event) {
        const resumeTopic = this.getCurrentConversation(event);
        this.conversationService.setRerunCurrentConversation(resumeTopic);
        this.isShowSidebar = false;
        this.blurInput();
    }
    resetColorAreas() {
        this.isOverDeleteConvArea = false;
        this.isOverCurrentConvArea = false;
    }
};
__decorate([
    core_1.ViewChild('search'),
    __metadata("design:type", core_1.ElementRef)
], UnifinishedTopicComponent.prototype, "inputEl", void 0);
UnifinishedTopicComponent = __decorate([
    core_1.Component({
        selector: 'oty-unfinished-topic',
        template: __webpack_require__("./src/app/ontology.unfinished.topic/unfinished.topic.component.html"),
        styles: [__webpack_require__("./src/app/ontology.unfinished.topic/unfinished.topic.component.scss")]
    }),
    __metadata("design:paramtypes", [conversation_service_1.ConversationService])
], UnifinishedTopicComponent);
exports.UnifinishedTopicComponent = UnifinishedTopicComponent;


/***/ }),

/***/ "./src/app/ontology.unfinished.topic/unfinished.topic.module.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const common_1 = __webpack_require__("./node_modules/@angular/common/esm2015/common.js");
const shared_module_1 = __webpack_require__("./src/app/common/shared.module.ts");
const primeng_1 = __webpack_require__("./node_modules/primeng/primeng.js");
const unfinished_topic_component_1 = __webpack_require__("./src/app/ontology.unfinished.topic/unfinished.topic.component.ts");
let UnfinishedTopicModule = class UnfinishedTopicModule {
};
UnfinishedTopicModule = __decorate([
    core_1.NgModule({
        imports: [
            common_1.CommonModule,
            shared_module_1.SharedModule,
            primeng_1.DragDropModule
        ],
        exports: [unfinished_topic_component_1.UnifinishedTopicComponent],
        declarations: [unfinished_topic_component_1.UnifinishedTopicComponent]
    })
], UnfinishedTopicModule);
exports.UnfinishedTopicModule = UnfinishedTopicModule;


/***/ }),

/***/ "./src/app/services/auth-guard.service.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const auth_service_1 = __webpack_require__("./src/app/services/auth.service.ts");
let AuthGuardService = class AuthGuardService {
    constructor(authService, router) {
        this.authService = authService;
        this.router = router;
    }
    canLoad(route) {
        if (this.authService.isAuthenticated()) {
            return true;
        }
        else {
            this.router.navigate(['/login']);
            return false;
        }
    }
    canActivate(route, state) {
        if (this.authService.isAuthenticated()) {
            return true;
        }
        else {
            this.router.navigate(['/login']);
            return false;
        }
    }
};
AuthGuardService = __decorate([
    core_1.Injectable(),
    __metadata("design:paramtypes", [auth_service_1.AuthService,
        router_1.Router])
], AuthGuardService);
exports.AuthGuardService = AuthGuardService;


/***/ }),

/***/ "./src/app/services/auth.service.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const global = __webpack_require__("./src/app/common/global.ts");
const http_1 = __webpack_require__("./node_modules/@angular/common/esm2015/http.js");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
let AuthService = class AuthService {
    constructor(http, router) {
        this.http = http;
        this.router = router;
        this.usernameEvent = new core_1.EventEmitter();
    }
    login(username, password) {
        const body = new http_1.HttpParams()
            .set('username', username)
            .set('password', password);
        return this.http.post(global.BASE_API_URL + '/convbot/users/login', body.toString(), {
            headers: new http_1.HttpHeaders({
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authorization': "Basic " + btoa(username + ':' + password)
            }),
            observe: 'response'
        }).map(res => {
            let token = res.headers.get('Authorization');
            if (token) {
                this.saveSession(username, token);
            }
            return token;
        });
    }
    isAuthenticated() {
        return sessionStorage.getItem('otyToken') !== null;
    }
    getUsername() {
        return sessionStorage.getItem('otyUsername');
    }
    logout() {
        return this.http.get(global.BASE_API_URL + '/convbot/users/logout').map((res) => {
            this.deleteSession();
            return res.success;
        });
    }
    saveSession(username, token) {
        sessionStorage.setItem('otyUsername', username);
        sessionStorage.setItem('otyToken', token);
        this.usernameEvent.emit(username);
    }
    deleteSession() {
        sessionStorage.removeItem('otyUsername');
        sessionStorage.removeItem('otyToken');
        this.usernameEvent.emit(null);
    }
};
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], AuthService.prototype, "usernameEvent", void 0);
AuthService = __decorate([
    core_1.Injectable(),
    __metadata("design:paramtypes", [http_1.HttpClient, router_1.Router])
], AuthService);
exports.AuthService = AuthService;


/***/ }),

/***/ "./src/app/services/conversation.service.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const http_1 = __webpack_require__("./node_modules/@angular/common/esm2015/http.js");
const global = __webpack_require__("./src/app/common/global.ts");
let ConversationService = class ConversationService {
    constructor(http) {
        this.http = http;
        this.onHistoryConversationChange = new core_1.EventEmitter();
        this.onRerunCurrentConversation = new core_1.EventEmitter();
        this.onFormatConversation = new core_1.EventEmitter();
        this.scrollConversationHistoryList = new core_1.EventEmitter();
        this.conversationHistoryList = [];
    }
    startConversation(projectId) {
        return this.http.get(global.BASE_API_URL + '/convbot/conversation/' + projectId).map((res) => res.data);
    }
    historyConversationChange() {
        this.onHistoryConversationChange.emit();
    }
    setRerunCurrentConversation(unfinishedTopic) {
        this.onRerunCurrentConversation.emit(unfinishedTopic);
    }
    setDialogHistory(dialogHistory) {
        this.dialogHistory = dialogHistory;
    }
    getDialogHistory() {
        return this.dialogHistory;
    }
    continueConversation(projectId, userRequest, userControlType) {
        let dialogId = this.getDialogHistory().slice(-1)[0].dialogId;
        userControlType = userControlType || "conversation";
        let conversation = {
            dialogHistory: this.getDialogHistory(),
            userControlType: userControlType,
            userRequest: userRequest
        };
        return this.http.post(global.BASE_API_URL + '/convbot/conversation/' + projectId + "/" + dialogId, conversation)
            .map((res) => res.data);
    }
    getGraphData(projectId) {
        return this.http.get(global.BASE_API_URL + '/convbot/graph/' + projectId)
            .map((res) => res.data);
    }
    saveGraphData(projectId, graphData) {
        return this.http.post(global.BASE_API_URL + '/convbot/graph/' + projectId, graphData)
            .map((res) => res.data);
    }
    initConversationHistoryList(conversations) {
        this.conversationHistoryList = conversations;
    }
    addConversationHistoryList(res) {
        let lastDialogHistory = res.dialogHistory.slice(-1)[0];
        this.latestConversationData = res;
        this.conversationHistoryList.push({
            dialogState: lastDialogHistory,
            content: {
                botGreetingCaption: res.botGreetingCaption,
                botResponseCaption: res.botResponseCaption
            }
        });
        this.scrollConversationHistoryList.emit();
    }
    getConversationHistoryList() {
        return this.conversationHistoryList;
    }
    clearConversationHistoryList() {
        this.conversationHistoryList = [];
    }
    getUnfinishedTopics() {
        let arr;
        if (this.latestConversationData && this.latestConversationData.botResponseData.value.unfinishedTopics) {
            arr = this.latestConversationData.botResponseData.value.unfinishedTopics;
        }
        return arr;
    }
    fireEventFormatConversation(res) {
        this.onFormatConversation.emit(res);
    }
    /**
     * rerun a conversation from restored conversation list
     * @param {string} projectId
     * @param {string} topicId
     * @returns {Observable<ConversationRes>}
     */
    resumeTopic(projectId, topicId) {
        let userRequest = {
            "key": "command",
            "type": "complex",
            "value": {
                "name": "resumeTopic",
                "component": {
                    "topicId": topicId
                }
            }
        };
        return this.continueConversation(projectId, userRequest, "uiEvent");
    }
    getEntityDetail(projectId, entityName) {
        let userRequest = {
            "key": "command",
            "type": "complex",
            "value": {
                "name": "graphAction",
                "component": {
                    "type": "entity",
                    "entityName": entityName
                }
            }
        };
        return this.continueConversation(projectId, userRequest, "uiEvent");
    }
    deleteEntity(projectId, entityName) {
        let userRequest = {
            "key": "command",
            "type": "complex",
            "value": {
                "name": "deleteAction",
                "component": {
                    "type": "entity",
                    "entityName": entityName
                }
            }
        };
        return this.continueConversation(projectId, userRequest, "uiEvent");
    }
    //**
    getAttributeDetail(projectId, entityName, attributeId) {
        let userRequest = {
            "key": "command",
            "type": "complex",
            "value": {
                "name": "graphAction",
                "component": {
                    "type": "attribute",
                    "entityName": entityName,
                    "attributeId": attributeId
                }
            }
        };
        return this.continueConversation(projectId, userRequest, "uiEvent");
    }
    /*
     getAttributeDetail(projectId: string, entityName: string, attributeName:string): Observable<ConversationRes> {
      let userRequest: UserRequest = {
        "key" : "command",
        "type": "complex",
        "value": {
          "name": "graphAction",
          "component" : {
            "type" : "attribute",
            "entityName" : entityName,
            "attributeName": attributeName
          }
        }
      };
      return this.continueConversation(projectId, userRequest, "uiEvent");
    }
     */
    deleteAttribute(projectId, entityName, attributeName) {
        let userRequest = {
            "key": "command",
            "type": "complex",
            "value": {
                "name": "deleteAction",
                "component": {
                    "type": "attribute",
                    "entityName": entityName,
                    "attributeName": attributeName
                }
            }
        };
        return this.continueConversation(projectId, userRequest, "uiEvent");
    }
    getRelationDetail(projectId, relationName, sourceEntityName, targetEntityName) {
        let userRequest = {
            "key": "command",
            "type": "complex",
            "value": {
                "name": "graphAction",
                "component": {
                    "type": "relation",
                    "relationName": relationName,
                    "sourceEntityName": sourceEntityName,
                    "targetEntityName": targetEntityName
                }
            }
        };
        return this.continueConversation(projectId, userRequest, "uiEvent");
    }
    deleteRelation(projectId, relationName, sourceEntityName, targetEntityName) {
        let userRequest = {
            "key": "command",
            "type": "complex",
            "value": {
                "name": "deleteAction",
                "component": {
                    "type": "relation",
                    "relationName": relationName,
                    "sourceEntityName": sourceEntityName,
                    "targetEntityName": targetEntityName
                }
            }
        };
        return this.continueConversation(projectId, userRequest, "uiEvent");
    }
    deleteDefinition(projectId, entityName, index) {
        let userRequest = {
            "key": "command",
            "type": "complex",
            "value": {
                "name": "deleteAction",
                "component": {
                    "type": "definition",
                    "entityName": entityName,
                    "definition": index,
                }
            }
        };
        return this.continueConversation(projectId, userRequest, "uiEvent");
    }
};
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], ConversationService.prototype, "onHistoryConversationChange", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], ConversationService.prototype, "onRerunCurrentConversation", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], ConversationService.prototype, "onFormatConversation", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], ConversationService.prototype, "scrollConversationHistoryList", void 0);
ConversationService = __decorate([
    core_1.Injectable(),
    __metadata("design:paramtypes", [http_1.HttpClient])
], ConversationService);
exports.ConversationService = ConversationService;


/***/ }),

/***/ "./src/app/services/cytoscape.service.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const http_1 = __webpack_require__("./node_modules/@angular/common/esm2015/http.js");
const $ = __webpack_require__("./node_modules/jquery/dist/jquery.js");
const cytoscape = __webpack_require__("./node_modules/cytoscape/dist/cytoscape.js");
const coseBilkent = __webpack_require__("./node_modules/cytoscape-cose-bilkent/cytoscape-cose-bilkent.js");
const expandCollapse = __webpack_require__("./node_modules/cytoscape-expand-collapse/cytoscape-expand-collapse.js");
const global_local = __webpack_require__("./src/app/common/global_local.ts");
const global_remote = __webpack_require__("./src/app/common/global_remote.ts");
__webpack_require__("./node_modules/rxjs/_esm2015/add/operator/map.js");
let CytoscapeService = class CytoscapeService {
    // Receive a div element. This will be the drawing space
    constructor(http) {
        this.http = http;
        this.event = new core_1.EventEmitter();
        this.elementSelected = null;
        this.sourceType = null;
        //event names match interface ./app/services/cytoscape.interface.ts
        this.eventName = {
            "ONNOTHINGCLICK": "onNothingClick",
            "ONNODECLICK": "onNodeClick",
            "ONATTRIBUTECLICK": "onAttributeClick",
            "ONEDGECLICK": "onEdgeClick",
            "ONEDGECREATE": "onEdgeCreate"
        };
    }
    initCytoscape(cyElement) {
        this.cy = cytoscape({
            container: cyElement
        });
        /*
        if (!this.cy.edgehandles) {
          cytoscape.use(edgehandles);
          this.cy.edgehandles();
        }
        */
        if (typeof cytoscape('core', 'expandCollapse') !== 'function') {
            cytoscape.use(expandCollapse, $);
        }
        cytoscape.use(coseBilkent);
        this.loadStyleJson();
        this.loadExpandCollapse();
        this.bindEvents();
        return this.cy;
    }
    loadExpandCollapse() {
        this.cy.expandCollapse({
            layoutBy: {
                name: 'cose-bilkent',
                animate: true,
                randomize: false,
                fit: false,
                nodeDimensionsIncludeLabels: true,
                tilingPaddingVertical: 50,
                tilingPaddingHorizontal: 50,
            },
            animate: true,
            undoable: false
        });
    }
    loadStyleJson() {
        this.cy.style(global_local.cytoscape_style);
    }
    updateLayout(layoutName) {
        const options = {
            name: layoutName,
            fit: false,
            minNodeSpacing: 50,
            avoidOverlap: true,
            animate: false,
            stop: () => {
                this.cy.reset();
                this.cy.center();
                if (layoutName === 'breadthfirst') {
                    this.cy.fit();
                }
            }
        };
        if (layoutName === 'cose-bilkent') {
            options.nodeDimensionsIncludeLabels = true;
            options.tilingPaddingVertical = 50;
            options.tilingPaddingHorizontal = 50;
        }
        const layout = this.cy.layout(options);
        layout.run();
    }
    resize() {
        this.cy.resize();
    }
    run(layoutName) {
        console.log('[cytoscape.service::run]--sourceType=' + this.sourceType);
        this.cy.layout({
            name: layoutName || 'cose-bilkent',
            stop: () => {
                this.cy.reset();
                this.cy.center();
            }
        }).run();
    }
    bindEvents() {
        console.log('[cytoscape.service::center]--sourceType=' + this.sourceType);
        this.bindClickEvent();
        this.bindEdgeEvent();
        this.bindMouseOverEvent();
        this.bindMouseOutEvent();
    }
    zoom(zoomNumber) {
        //console.log('[cytoscape.service::zoom]--sourceType='+this.sourceType);
        const pos = this.cy.nodes().position();
        this.cy.zoom({
            level: zoomNumber + this.cy.zoom(),
            position: pos
        });
    }
    center() {
        console.log('[cytoscape.service::center]--sourceType=' + this.sourceType);
        this.cy.reset();
        this.cy.center();
    }
    bindClickEvent() {
        console.log('[cytoscape.service::bindClickEvent]--sourceType=' + this.sourceType);
        this.cy.on('click select', (evt) => {
            let eventName = "", target = evt.target;
            if (target.length === 1) {
                if (evt.type === "select" || (evt.type === "click" && target.selected())) {
                    this.resetNodesStyle();
                    // clicked Node: add class for highlight(click), highlight current node, highlight related edge.
                    if (target.isNode()) {
                        this.elementSelected = target;
                        target.addClass('node-on-select-highlight');
                        target.connectedEdges().addClass('node-on-select-edge-highlight');
                        if (target.isChild()) {
                            eventName = this.eventName.ONATTRIBUTECLICK;
                        }
                        else {
                            eventName = this.eventName.ONNODECLICK;
                        }
                    }
                    // clicked Edge: add class for highlight(click), highlight current edge, highlight related node
                    if (target.isEdge()) {
                        this.elementSelected = target;
                        target.addClass('node-on-select-edge-highlight');
                        target.connectedNodes().addClass('node-on-select-highlight');
                        eventName = this.eventName.ONEDGECLICK;
                    }
                }
            }
            else {
                if (evt.type === "click") {
                    this.resetNodesStyle();
                    this.elementSelected = null;
                    eventName = this.eventName.ONNOTHINGCLICK;
                }
            }
            if (eventName) {
                this.event.emit({
                    name: eventName,
                    event: evt
                });
            }
        });
    }
    bindEdgeEvent() {
        console.log('[cytoscape.service::bindEdgeEvent]--sourceType=' + this.sourceType);
        this.cy.on('ehcomplete', (event, sourceNode, targetNode, addedEles) => {
            this.event.emit({
                name: this.eventName.ONEDGECREATE,
                event: event,
                data: {
                    sourceNode: sourceNode,
                    targetNode: targetNode,
                    addedEles: addedEles
                }
            });
        });
    }
    resetNodesStyle() {
        console.log('[cytoscape.service::resetNodesStyle]--sourceType=' + this.sourceType);
        const nodes = this.cy.nodes();
        nodes.removeClass('node-on-select-highlight');
        nodes.connectedEdges().removeClass('node-on-select-edge-highlight');
        // Needed or else the on-hover classes remains after the on-click event occurs
        nodes.removeClass('node-on-hover-highlight');
        nodes.connectedEdges().removeClass('node-on-hover-edge-highlight');
    }
    // mouse over: add class for highlight(click), highlight current edge, highlight related node
    bindMouseOverEvent() {
        this.cy.on('mouseover', ':selectable', (evt) => {
            const target = evt.target;
            target.addClass('node-on-hover-highlight');
            target.addClass('node-on-hover-edge-highlight');
            target.connectedNodes().addClass('node-on-hover-highlight');
            target.connectedEdges().addClass('node-on-hover-edge-highlight');
        });
    }
    // mouse out: remove class for highlight(click), highlight current edge, highlight related node.
    bindMouseOutEvent() {
        console.log('[cytoscape.service::bindMouseOutEvent]--sourceType=' + this.sourceType);
        this.cy.on('mouseout', ':selectable', evt => {
            const target = evt.target;
            target.removeClass('node-on-hover-highlight');
            target.removeClass('node-on-hover-edge-highlight');
            target.connectedNodes().removeClass('node-on-hover-highlight');
            target.connectedEdges().removeClass('node-on-hover-edge-highlight');
        });
    }
    fireDeleteEvent() {
        this.event.emit({ name: 'onDelete', event: null });
    }
    deleteElement(id) {
        this.cy.$(`#${id}`).remove();
    }
    selectElement(id) {
        this.resetNodesStyle();
        this.cy.$(`#${id}`).select();
    }
    collapseAll() {
        this.cy.expandCollapse('get').collapseAll();
    }
    expandAll() {
        //console.log('[cytoscape.service::expandAll]--sourceType='+this.sourceType);
        this.redrawGraph(this.sourceType);
        //this.cy.expandCollapse('get').expandAll();
    }
    redrawGraph(sourceType) {
        console.log('[cytoscape.service::redrawGroup]--sourceType=' + this.sourceType);
        let json = this.cy.json();
        if ("remote" === this.sourceType) {
            json = this.cy.style(global_remote.cytoscape_style);
        }
        else {
            json = this.cy.style(global_local.cytoscape_style);
        }
        this.cy.elements().remove();
        this.cy.json(json);
    }
};
CytoscapeService = __decorate([
    core_1.Injectable(),
    __metadata("design:paramtypes", [http_1.HttpClient])
], CytoscapeService);
exports.CytoscapeService = CytoscapeService;


/***/ }),

/***/ "./src/app/services/notification.service.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
let NotificationService = class NotificationService {
    constructor() {
        this.notificationEvent = new core_1.EventEmitter();
    }
    showNotification(type, message) {
        this.notificationEvent.emit({ 'type': type, 'message': message });
    }
};
NotificationService = __decorate([
    core_1.Injectable()
], NotificationService);
exports.NotificationService = NotificationService;


/***/ }),

/***/ "./src/app/services/project.service.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const http_1 = __webpack_require__("./node_modules/@angular/common/esm2015/http.js");
const global = __webpack_require__("./src/app/common/global.ts");
let ProjectService = class ProjectService {
    constructor(http) {
        this.http = http;
    }
    getProjects() {
        return this.http.get(global.BASE_API_URL + '/convbot/projects')
            .map((res) => res.data);
    }
    getProject(id) {
        return this.http.get(global.BASE_API_URL + '/convbot/project/' + id)
            .map((res) => {
            return res.data;
        });
    }
    createProject(name) {
        let reqData = {
            'projectName': name
        };
        return this.http.put(global.BASE_API_URL + '/convbot/project', reqData);
    }
    updateProject(project) {
        return this.http.post(global.BASE_API_URL + '/convbot/project/' + project.projectId, project)
            .map((res) => {
            return res.data;
        });
    }
    ;
    updateProjectName(id, name) {
        let project = {
            'projectName': name
        };
        return this.http.post(global.BASE_API_URL + '/convbot/project/' + id, project)
            .map((res) => {
            return res.data;
        });
    }
    deleteProject(id) {
        return this.http.delete(global.BASE_API_URL + '/convbot/project/' + id);
    }
    exportProject(id) {
        return this.http.get(global.BASE_API_URL + '/convbot/project/export/' + id);
    }
};
ProjectService = __decorate([
    core_1.Injectable(),
    __metadata("design:paramtypes", [http_1.HttpClient])
], ProjectService);
exports.ProjectService = ProjectService;


/***/ }),

/***/ "./src/environments/environment.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.
Object.defineProperty(exports, "__esModule", { value: true });
exports.environment = {
    production: false
};


/***/ }),

/***/ "./src/main.ts":
/***/ (function(module, exports, __webpack_require__) {

"use strict";

Object.defineProperty(exports, "__esModule", { value: true });
const core_1 = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
const platform_browser_dynamic_1 = __webpack_require__("./node_modules/@angular/platform-browser-dynamic/esm2015/platform-browser-dynamic.js");
const app_module_1 = __webpack_require__("./src/app/app.module.ts");
const environment_1 = __webpack_require__("./src/environments/environment.ts");
if (environment_1.environment.production) {
    core_1.enableProdMode();
}
platform_browser_dynamic_1.platformBrowserDynamic().bootstrapModule(app_module_1.AppModule);


/***/ }),

/***/ 0:
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__("./src/main.ts");


/***/ })

},[0]);
//# sourceMappingURL=main.bundle.js.map