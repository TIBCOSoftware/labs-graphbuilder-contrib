webpackJsonp(["projects.module"],{

/***/ "./src/app/ontology.projects/projects-routing.module.ts":
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
const projects_component_1 = __webpack_require__("./src/app/ontology.projects/projects.component.ts");
const auth_guard_service_1 = __webpack_require__("./src/app/services/auth-guard.service.ts");
let ProjectsRoutingModule = class ProjectsRoutingModule {
};
ProjectsRoutingModule = __decorate([
    core_1.NgModule({
        imports: [router_1.RouterModule.forChild([
                { path: '', component: projects_component_1.ProjectsComponent, canActivate: [auth_guard_service_1.AuthGuardService] },
                { path: 'project/:id', loadChildren: 'app/ontology.project/project.module#ProjectModule' },
            ])],
        exports: [router_1.RouterModule]
    })
], ProjectsRoutingModule);
exports.ProjectsRoutingModule = ProjectsRoutingModule;


/***/ }),

/***/ "./src/app/ontology.projects/projects.component.html":
/***/ (function(module, exports) {

module.exports = "<div id=\"oty-projects-bar\">\n  <div id=\"oty-projects-title\" class=\"pull-left\">\n    <p>Projects</p>\n    <div><span>{{projectList?.length || 0}}</span>items displayed</div>\n  </div>\n  <div id=\"oty-projects-toolbar\" class=\"pull-right\">\n    <div class=\"tc-search pull-left\">\n      <div class=\"tc-text-areas tc-search-container\">\n        <input #searchInput type='text' class=\"tc-search inputIcon\" placeholder=\"Search...\"\n               (input)=\"onSearchChange($event.target.value)\"/>\n        <span class=\"tc-search-icon\"></span>\n      </div>\n    </div>\n\n    <button routerLink='/project-new/' class=\"tc-buttons tc-button-icon-text\">\n      <i class=\"oty-icon-add\"></i>\n      <span class=\"tc-button-icon-text-title\">Create New project</span>\n    </button>\n  </div>\n</div>\n\n<div id='oty-projects-list'>\n  <ul class=\"list-unstyled\">\n    <li *ngFor='let project of projectList; let i = index' [ngClass]=\"{'oty-projects-list-edge': (i+1) % 4 === 0}\"\n        (onDragStart)=\"onDragStart($event, project)\" (onDragEnd)=\"onDragEnd($event)\" pDraggable=\"oty-project\">\n      <a routerLink='/project/{{ project.projectId }}'>\n        <img src=\"../../assets/images/project/project{{i+1}}.svg\"\n             #img\n             (error)=\"img.src='../../assets/images/project/project'+ getRandomProjectIcon() +'.svg'\"/>\n        <span>{{ project.projectName }}</span>\n      </a>\n    </li>\n    <li>\n      <a class='oty-add-btn' routerLink='/project-new'>\n        <i class=\"oty-icon-add\"></i>\n      </a>\n    </li>\n    <li pDroppable=\"oty-project\" *ngIf=\"isShowDeleteArea\"\n        (onDrop)=\"deleteProject($event)\">\n      <a class='oty-add-btn oty-remove-project' routerLink='/project-new' [ngClass]=\"{'delete-area-selected':isOverDelete}\">\n        <i pDroppable=\"oty-project\"\n           (onDragEnter)=\"isOverDelete = true\"\n           (onDragLeave)=\"isOverDelete = false\" class=\"oty-icon-zoom-out\"></i>\n      </a>\n    </li>\n  </ul>\n  <p *ngIf=\"!projectList\">No project.</p>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.projects/projects.component.scss":
/***/ (function(module, exports) {

module.exports = ".display-modal {\n  display: block;\n  left: 0; }\n  .display-modal .modal-content {\n    left: 0;\n    top: 35%; }\n  :host {\n  width: 978px;\n  height: 100%;\n  margin: 0 auto;\n  display: block;\n  padding-top: 42px; }\n  :host #oty-projects-bar {\n    height: 50px;\n    margin-bottom: 36px; }\n  :host #oty-projects-bar #oty-projects-title p {\n      float: left;\n      display: inline-block;\n      margin: 0;\n      font-size: 30px;\n      font-weight: 600; }\n  :host #oty-projects-bar #oty-projects-title > div {\n      display: inline-block;\n      line-height: 50px; }\n  :host #oty-projects-bar #oty-projects-title > div span {\n        margin: 0 10px 0 20px;\n        border-radius: 100px;\n        padding: 3px 10px;\n        background: #DEDEDE;\n        color: #212121; }\n  :host #oty-projects-bar #oty-projects-toolbar .tc-search {\n      margin: 7px 20px 0 0; }\n  :host #oty-projects-bar #oty-projects-toolbar .tc-search .tc-search {\n        padding-left: 30px;\n        background: transparent; }\n  :host #oty-projects-bar #oty-projects-toolbar .tc-search .tc-text-areas {\n        border-bottom: 2px #DEDEDE solid !important; }\n  :host #oty-projects-bar #oty-projects-toolbar .tc-search .tc-search-icon {\n        left: 0;\n        right: auto;\n        top: 7px; }\n  :host #oty-projects-bar #oty-projects-toolbar .tc-button-icon-text i {\n      margin-right: 8px; }\n  :host #oty-projects-list {\n    clear: both; }\n  :host #oty-projects-list ul li {\n      margin: 0 24px 24px 0;\n      float: left;\n      display: inline-block; }\n  :host #oty-projects-list ul li.oty-projects-list-edge {\n        margin-right: 0; }\n  :host #oty-projects-list ul li a {\n        height: 194px;\n        width: 226px;\n        display: block;\n        border: 1px solid #DEDEDE;\n        border-radius: 4px;\n        -webkit-box-shadow: 0 2px 4px 0 #B6B6B6;\n                box-shadow: 0 2px 4px 0 #B6B6B6;\n        text-align: center; }\n  :host #oty-projects-list ul li a img {\n          display: block;\n          margin: 20px auto; }\n  :host #oty-projects-list ul li a span {\n          height: 36px;\n          width: 106px;\n          font-size: 18px;\n          letter-spacing: 0.3px;\n          line-height: 36px; }\n  :host #oty-projects-list ul li .oty-add-btn {\n        text-decoration: none;\n        border: 2px dashed #DEDEDE;\n        border-radius: 3px;\n        -webkit-box-shadow: none;\n                box-shadow: none; }\n  :host #oty-projects-list ul li .oty-add-btn i {\n          text-align: center;\n          font-size: 56px;\n          color: #727272;\n          position: relative;\n          top: calc(40% - 10px); }\n  :host #oty-projects-list ul li .oty-remove-project {\n        border-color: #FF0000; }\n  :host #oty-projects-list ul li .oty-remove-project i {\n          color: #D0021B;\n          padding: 65px 80px 65px 80px; }\n  :host #oty-projects-list ul li .delete-area-selected {\n        background: #FFB6BF; }\n  :host #oty-projects-list p {\n      height: 500px;\n      line-height: 500px;\n      font-size: 36px;\n      text-align: center; }\n"

/***/ }),

/***/ "./src/app/ontology.projects/projects.component.ts":
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
const project_service_1 = __webpack_require__("./src/app/services/project.service.ts");
const notification_service_1 = __webpack_require__("./src/app/services/notification.service.ts");
let ProjectsComponent = class ProjectsComponent {
    constructor(projectService, notificationService) {
        this.projectService = projectService;
        this.notificationService = notificationService;
        this.isShowDeleteArea = false;
        this.isOverDelete = false;
    }
    ngOnInit() {
        this.loadProjects();
    }
    loadProjects() {
        this.projectService.getProjects()
            .subscribe((projects) => {
            this.projectList = projects;
            this.projectListOriginal = projects;
        });
    }
    onSearchChange(value) {
        this.projectList = this.projectListOriginal.filter((item) => {
            return item.projectName.toLowerCase().includes(value.toLowerCase());
        });
    }
    getRandomProjectIcon() {
        let max = 6, n = Math.floor(Math.random() * Math.floor(max)) + 1;
        return n;
    }
    onDragStart(event, project) {
        this.isShowDeleteArea = true;
        event.dataTransfer.setData('projectSelected', JSON.stringify(project));
    }
    onDragEnd(event) {
        this.isShowDeleteArea = false;
    }
    deleteProject(event) {
        const project = JSON.parse(event.dataTransfer.getData('projectSelected'));
        this.projectService
            .deleteProject(project.projectId)
            .subscribe((res) => {
            if (res.success === true) {
                this.notificationService.showNotification('success', `${project.projectName} deleted`);
            }
            this.loadProjects();
            this.clearSearchInput();
            this.isOverDelete = false;
        });
    }
    clearSearchInput() {
        this.searchInput.nativeElement.value = '';
    }
};
__decorate([
    core_1.ViewChild('searchInput'),
    __metadata("design:type", core_1.ElementRef)
], ProjectsComponent.prototype, "searchInput", void 0);
ProjectsComponent = __decorate([
    core_1.Component({
        selector: 'oty-projects',
        template: __webpack_require__("./src/app/ontology.projects/projects.component.html"),
        styles: [__webpack_require__("./src/app/ontology.projects/projects.component.scss")]
    }),
    __metadata("design:paramtypes", [project_service_1.ProjectService,
        notification_service_1.NotificationService])
], ProjectsComponent);
exports.ProjectsComponent = ProjectsComponent;


/***/ }),

/***/ "./src/app/ontology.projects/projects.module.ts":
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
const projects_component_1 = __webpack_require__("./src/app/ontology.projects/projects.component.ts");
const projects_routing_module_1 = __webpack_require__("./src/app/ontology.projects/projects-routing.module.ts");
const primeng_1 = __webpack_require__("./node_modules/primeng/primeng.js");
let ProjectsModule = class ProjectsModule {
};
ProjectsModule = __decorate([
    core_1.NgModule({
        imports: [common_1.CommonModule, projects_routing_module_1.ProjectsRoutingModule, primeng_1.DragDropModule],
        declarations: [projects_component_1.ProjectsComponent]
    })
], ProjectsModule);
exports.ProjectsModule = ProjectsModule;


/***/ })

});
//# sourceMappingURL=projects.module.chunk.js.map