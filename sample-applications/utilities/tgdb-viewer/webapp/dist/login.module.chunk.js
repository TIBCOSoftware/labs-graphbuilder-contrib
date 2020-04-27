webpackJsonp(["login.module"],{

/***/ "./src/app/ontology.login/login-routing.module.ts":
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
const login_component_1 = __webpack_require__("./src/app/ontology.login/login.component.ts");
let LoginRoutingModule = class LoginRoutingModule {
};
LoginRoutingModule = __decorate([
    core_1.NgModule({
        imports: [router_1.RouterModule.forChild([
                { path: '', component: login_component_1.LoginComponent }
            ])],
        exports: [router_1.RouterModule]
    })
], LoginRoutingModule);
exports.LoginRoutingModule = LoginRoutingModule;


/***/ }),

/***/ "./src/app/ontology.login/login.component.html":
/***/ (function(module, exports) {

module.exports = "<div class=\"oty-login\">\n  <div>\n  </div>\n  <div>\n    <span>Log in</span>\n    <form [formGroup]=\"elemForm\" (ngSubmit)=\"sendData(elemForm)\">\n      <div>\n        <input id=\"username\" type=\"text\" formControlName=\"username\" autofocus required/>\n        <label for=\"username\">User name</label>\n      </div>\n      <div>\n        <input id=\"password\" type=\"password\" formControlName=\"password\" required/>\n        <label for=\"password\">Password</label>\n      </div>\n      <div class=\"error-message\" *ngIf=\"errorMessage\">\n        {{errorMessage}}\n      </div>\n      <button type=\"submit\" [disabled]=\"!elemForm.valid\" [ngClass]=\"{'disabled': !elemForm.valid}\">\n        Log in\n      </button>\n    </form>\n  </div>\n</div>\n\n"

/***/ }),

/***/ "./src/app/ontology.login/login.component.scss":
/***/ (function(module, exports) {

module.exports = ".oty-login {\n  height: 100%;\n  left: 0;\n  position: fixed;\n  top: 0;\n  width: 100%;\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex; }\n  .oty-login > div:nth-child(1) {\n    background-color: #19414e;\n    height: 100%;\n    width: 40%;\n    background-image: url('tibco-labs-tgdb-terminal.5b944cfc0a9eb2ea7fa2.png');\n    background-repeat: no-repeat;\n    background-attachment: fixed;\n    background-position: 10% 227px; }\n  .oty-login > div:nth-child(2) {\n    height: 100%;\n    width: 60%; }\n  .oty-login > div:nth-child(2) span {\n      font-size: 24px;\n      height: 31px;\n      left: 20%;\n      line-height: 31px;\n      position: relative;\n      top: 190px;\n      width: 61px; }\n  .oty-login > div:nth-child(2) form {\n      font-size: 13px;\n      height: 286px;\n      left: 20%;\n      line-height: 17px;\n      position: relative;\n      top: 227px;\n      width: 404px; }\n  .oty-login > div:nth-child(2) form div {\n        position: relative; }\n  .oty-login > div:nth-child(2) form div label {\n          color: #DEDEDE;\n          font-size: 16px;\n          position: absolute;\n          top: 22px;\n          left: 0; }\n  .oty-login > div:nth-child(2) form div input {\n          border: 0;\n          border-bottom: 1px #DEDEDE solid;\n          height: 45px;\n          margin-bottom: 50px;\n          width: 100%;\n          padding-top: 12px; }\n  .oty-login > div:nth-child(2) form div input:focus {\n            border-bottom: 1px #2694D3 solid;\n            outline: none; }\n  .oty-login > div:nth-child(2) form div input:focus ~ label {\n            color: #2694D3;\n            font-size: 13px;\n            font-weight: normal;\n            top: 0;\n            left: 0; }\n  .oty-login > div:nth-child(2) form div input:valid ~ label {\n            display: none; }\n  .oty-login > div:nth-child(2) form button {\n        background-color: #0081cb;\n        border-radius: 5px;\n        color: #FFFFFF;\n        font-size: 16px;\n        height: 42px;\n        width: 225px; }\n  .oty-login > div:nth-child(2) form .disabled {\n        background-color: rgba(0, 129, 203, 0.44); }\n  .oty-login > div .error-message {\n    color: #D0021B;\n    position: relative;\n    width: 100%;\n    height: 60px;\n    font-size: 16px;\n    text-align: left; }\n  .oty-login > div .disabled {\n    background-color: rgba(0, 129, 203, 0.44); }\n"

/***/ }),

/***/ "./src/app/ontology.login/login.component.ts":
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
const forms_1 = __webpack_require__("./node_modules/@angular/forms/esm2015/forms.js");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const auth_service_1 = __webpack_require__("./src/app/services/auth.service.ts");
const notification_service_1 = __webpack_require__("./src/app/services/notification.service.ts");
let LoginComponent = class LoginComponent {
    constructor(fb, notificationService, authService, router) {
        this.notificationService = notificationService;
        this.authService = authService;
        this.router = router;
        this.formBuilder = fb;
    }
    ngOnInit() {
        this.elemForm = this.formBuilder.group({
            'username': ["", [forms_1.Validators.required, forms_1.Validators.pattern('^[a-zA-Z0-9.-]+$')]],
            'password': ["", [forms_1.Validators.required, forms_1.Validators.pattern('^[a-zA-Z0-9.-]+$')]],
        });
        if (sessionStorage.getItem('otyToken')) {
            this.authService.deleteSession();
            this.router.navigate(['/login']);
        }
    }
    sendData(form) {
        this.authService.login(form.value.username, form.value.password).subscribe((token) => {
            if (token) {
                sessionStorage.setItem('otyToken', token);
                this.router.navigate(['/']);
            }
        }, (error) => {
            this.errorMessage = 'Username or password incorrect';
        });
    }
};
LoginComponent = __decorate([
    core_1.Component({
        selector: 'oty-login',
        template: __webpack_require__("./src/app/ontology.login/login.component.html"),
        styles: [__webpack_require__("./src/app/ontology.login/login.component.scss")]
    }),
    __param(0, core_1.Inject(forms_1.FormBuilder)),
    __metadata("design:paramtypes", [forms_1.FormBuilder,
        notification_service_1.NotificationService,
        auth_service_1.AuthService,
        router_1.Router])
], LoginComponent);
exports.LoginComponent = LoginComponent;


/***/ }),

/***/ "./src/app/ontology.login/login.module.ts":
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
const login_component_1 = __webpack_require__("./src/app/ontology.login/login.component.ts");
const shared_module_1 = __webpack_require__("./src/app/common/shared.module.ts");
const login_routing_module_1 = __webpack_require__("./src/app/ontology.login/login-routing.module.ts");
let LoginModule = class LoginModule {
};
LoginModule = __decorate([
    core_1.NgModule({
        imports: [
            common_1.CommonModule,
            shared_module_1.SharedModule,
            login_routing_module_1.LoginRoutingModule
        ],
        declarations: [login_component_1.LoginComponent]
    })
], LoginModule);
exports.LoginModule = LoginModule;


/***/ })

});
//# sourceMappingURL=login.module.chunk.js.map