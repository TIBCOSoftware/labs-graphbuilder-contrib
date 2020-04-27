webpackJsonp(["project.module"],{

/***/ "./node_modules/angular-split/esm2015/angular-split.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "AngularSplitModule", function() { return AngularSplitModule; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "SplitComponent", function() { return SplitComponent; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "SplitAreaDirective", function() { return SplitAreaDirective; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ɵa", function() { return SplitGutterDirective; });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__angular_core__ = __webpack_require__("./node_modules/@angular/core/esm2015/core.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__angular_common__ = __webpack_require__("./node_modules/@angular/common/esm2015/common.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_2_rxjs_Subject__ = __webpack_require__("./node_modules/rxjs/_esm2015/Subject.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_3_rxjs_add_operator_debounceTime__ = __webpack_require__("./node_modules/rxjs/_esm2015/add/operator/debounceTime.js");





/**
 * @fileoverview added by tsickle
 * @suppress {checkTypes} checked by tsc
 */
/**
 * angular-split
 *
 * Areas size are set in percentage of the split container.
 * Gutters size are set in pixels.
 *
 * So we set css 'flex-basis' property like this (where 0 <= area.size <= 1):
 *  calc( { area.size * 100 }% - { area.size * nbGutter * gutterSize }px );
 *
 * Examples with 3 visible areas and 2 gutters:
 *
 * |                     10px                   10px                                  |
 * |---------------------[]---------------------[]------------------------------------|
 * |  calc(20% - 4px)          calc(20% - 4px)              calc(60% - 12px)          |
 *
 *
 * |                          10px                        10px                        |
 * |--------------------------[]--------------------------[]--------------------------|
 * |  calc(33.33% - 6.667px)      calc(33.33% - 6.667px)      calc(33.33% - 6.667px)  |
 *
 *
 * |10px                                                  10px                        |
 * |[]----------------------------------------------------[]--------------------------|
 * |0                 calc(66.66% - 13.333px)                  calc(33%% - 6.667px)   |
 *
 *
 *  10px 10px                                                                         |
 * |[][]------------------------------------------------------------------------------|
 * |0 0                               calc(100% - 20px)                               |
 *
 */
class SplitComponent {
    /**
     * @param {?} ngZone
     * @param {?} elRef
     * @param {?} cdRef
     * @param {?} renderer
     */
    constructor(ngZone, elRef, cdRef, renderer) {
        this.ngZone = ngZone;
        this.elRef = elRef;
        this.cdRef = cdRef;
        this.renderer = renderer;
        this._direction = 'horizontal';
        this._useTransition = false;
        this._disabled = false;
        this._width = null;
        this._height = null;
        this._gutterSize = 11;
        this._gutterColor = '';
        this._gutterImageH = '';
        this._gutterImageV = '';
        this._dir = 'ltr';
        this.dragStart = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["EventEmitter"](false);
        this.dragProgress = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["EventEmitter"](false);
        this.dragEnd = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["EventEmitter"](false);
        this.gutterClick = new __WEBPACK_IMPORTED_MODULE_0__angular_core__["EventEmitter"](false);
        this.transitionEndInternal = new __WEBPACK_IMPORTED_MODULE_2_rxjs_Subject__["Subject"]();
        this.transitionEnd = (/** @type {?} */ (this.transitionEndInternal.asObservable())).debounceTime(20);
        this.isViewInitialized = false;
        this.isDragging = false;
        this.draggingWithoutMove = false;
        this.currentGutterNum = 0;
        this.displayedAreas = [];
        this.hidedAreas = [];
        this.dragListeners = [];
        this.dragStartValues = {
            sizePixelContainer: 0,
            sizePixelA: 0,
            sizePixelB: 0,
            sizePercentA: 0,
            sizePercentB: 0,
        };
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set direction(v) {
        v = (v === 'vertical') ? 'vertical' : 'horizontal';
        this._direction = v;
        [...this.displayedAreas, ...this.hidedAreas].forEach(area => {
            area.comp.setStyleVisibleAndDir(area.comp.visible, this.isDragging, this.direction);
        });
        this.build(false, false);
    }
    /**
     * @return {?}
     */
    get direction() {
        return this._direction;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set useTransition(v) {
        v = (typeof (v) === 'boolean') ? v : (v === 'false' ? false : true);
        this._useTransition = v;
    }
    /**
     * @return {?}
     */
    get useTransition() {
        return this._useTransition;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set disabled(v) {
        v = (typeof (v) === 'boolean') ? v : (v === 'false' ? false : true);
        this._disabled = v;
        // Force repaint if modified from TS class (instead of the template)
        this.cdRef.markForCheck();
    }
    /**
     * @return {?}
     */
    get disabled() {
        return this._disabled;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set width(v) {
        v = Number(v);
        this._width = (!isNaN(v) && v > 0) ? v : null;
        this.build(false, false);
    }
    /**
     * @return {?}
     */
    get width() {
        return this._width;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set height(v) {
        v = Number(v);
        this._height = (!isNaN(v) && v > 0) ? v : null;
        this.build(false, false);
    }
    /**
     * @return {?}
     */
    get height() {
        return this._height;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set gutterSize(v) {
        v = Number(v);
        this._gutterSize = (!isNaN(v) && v > 0) ? v : 11;
        this.build(false, false);
    }
    /**
     * @return {?}
     */
    get gutterSize() {
        return this._gutterSize;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set gutterColor(v) {
        this._gutterColor = (typeof v === 'string' && v !== '') ? v : '';
        // Force repaint if modified from TS class (instead of the template)
        this.cdRef.markForCheck();
    }
    /**
     * @return {?}
     */
    get gutterColor() {
        return this._gutterColor;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set gutterImageH(v) {
        this._gutterImageH = (typeof v === 'string' && v !== '') ? v : '';
        // Force repaint if modified from TS class (instead of the template)
        this.cdRef.markForCheck();
    }
    /**
     * @return {?}
     */
    get gutterImageH() {
        return this._gutterImageH;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set gutterImageV(v) {
        this._gutterImageV = (typeof v === 'string' && v !== '') ? v : '';
        // Force repaint if modified from TS class (instead of the template)
        this.cdRef.markForCheck();
    }
    /**
     * @return {?}
     */
    get gutterImageV() {
        return this._gutterImageV;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set dir(v) {
        v = (v === 'rtl') ? 'rtl' : 'ltr';
        this._dir = v;
    }
    /**
     * @return {?}
     */
    get dir() {
        return this._dir;
    }
    /**
     * @return {?}
     */
    get cssFlexdirection() {
        return (this.direction === 'horizontal') ? 'row' : 'column';
    }
    /**
     * @return {?}
     */
    get cssWidth() {
        return this.width ? `${this.width}px` : '100%';
    }
    /**
     * @return {?}
     */
    get cssHeight() {
        return this.height ? `${this.height}px` : '100%';
    }
    /**
     * @return {?}
     */
    get cssMinwidth() {
        return (this.direction === 'horizontal') ? `${this.getNbGutters() * this.gutterSize}px` : null;
    }
    /**
     * @return {?}
     */
    get cssMinheight() {
        return (this.direction === 'vertical') ? `${this.getNbGutters() * this.gutterSize}px` : null;
    }
    /**
     * @return {?}
     */
    ngAfterViewInit() {
        this.isViewInitialized = true;
    }
    /**
     * @return {?}
     */
    getNbGutters() {
        return this.displayedAreas.length - 1;
    }
    /**
     * @param {?} comp
     * @return {?}
     */
    addArea(comp) {
        const /** @type {?} */ newArea = {
            comp,
            order: 0,
            size: 0,
        };
        if (comp.visible === true) {
            this.displayedAreas.push(newArea);
        }
        else {
            this.hidedAreas.push(newArea);
        }
        comp.setStyleVisibleAndDir(comp.visible, this.isDragging, this.direction);
        this.build(true, true);
    }
    /**
     * @param {?} comp
     * @return {?}
     */
    removeArea(comp) {
        if (this.displayedAreas.some(a => a.comp === comp)) {
            const /** @type {?} */ area = /** @type {?} */ (this.displayedAreas.find(a => a.comp === comp));
            this.displayedAreas.splice(this.displayedAreas.indexOf(area), 1);
            this.build(true, true);
        }
        else if (this.hidedAreas.some(a => a.comp === comp)) {
            const /** @type {?} */ area = /** @type {?} */ (this.hidedAreas.find(a => a.comp === comp));
            this.hidedAreas.splice(this.hidedAreas.indexOf(area), 1);
        }
    }
    /**
     * @param {?} comp
     * @param {?} resetOrders
     * @param {?} resetSizes
     * @return {?}
     */
    updateArea(comp, resetOrders, resetSizes) {
        // Only refresh if area is displayed (No need to check inside 'hidedAreas')
        const /** @type {?} */ item = this.displayedAreas.find(a => a.comp === comp);
        if (item) {
            this.build(resetOrders, resetSizes);
        }
    }
    /**
     * @param {?} comp
     * @return {?}
     */
    showArea(comp) {
        const /** @type {?} */ area = this.hidedAreas.find(a => a.comp === comp);
        if (area) {
            comp.setStyleVisibleAndDir(comp.visible, this.isDragging, this.direction);
            const /** @type {?} */ areas = this.hidedAreas.splice(this.hidedAreas.indexOf(area), 1);
            this.displayedAreas.push(...areas);
            this.build(true, true);
        }
    }
    /**
     * @param {?} comp
     * @return {?}
     */
    hideArea(comp) {
        const /** @type {?} */ area = this.displayedAreas.find(a => a.comp === comp);
        if (area) {
            comp.setStyleVisibleAndDir(comp.visible, this.isDragging, this.direction);
            const /** @type {?} */ areas = this.displayedAreas.splice(this.displayedAreas.indexOf(area), 1);
            areas.forEach(area => {
                area.order = 0;
                area.size = 0;
            });
            this.hidedAreas.push(...areas);
            this.build(true, true);
        }
    }
    /**
     * @param {?} resetOrders
     * @param {?} resetSizes
     * @return {?}
     */
    build(resetOrders, resetSizes) {
        this.stopDragging();
        // ¤ AREAS ORDER
        if (resetOrders === true) {
            // If user provided 'order' for each area, use it to sort them.
            if (this.displayedAreas.every(a => a.comp.order !== null)) {
                this.displayedAreas.sort((a, b) => (/** @type {?} */ (a.comp.order)) - (/** @type {?} */ (b.comp.order)));
            }
            // Then set real order with multiples of 2, numbers between will be used by gutters.
            this.displayedAreas.forEach((area, i) => {
                area.order = i * 2;
                area.comp.setStyleOrder(area.order);
            });
        }
        // ¤ AREAS SIZE PERCENT
        if (resetSizes === true) {
            const /** @type {?} */ totalUserSize = /** @type {?} */ (this.displayedAreas.reduce((total, s) => s.comp.size ? total + s.comp.size : total, 0));
            // If user provided 'size' for each area and total == 1, use it.
            if (this.displayedAreas.every(a => a.comp.size !== null) && totalUserSize > .999 && totalUserSize < 1.001) {
                this.displayedAreas.forEach(area => {
                    area.size = /** @type {?} */ (area.comp.size);
                });
            }
            else {
                const /** @type {?} */ size = 1 / this.displayedAreas.length;
                this.displayedAreas.forEach(area => {
                    area.size = size;
                });
            }
        }
        // ¤
        // If some real area sizes are less than gutterSize,
        // set them to zero and dispatch size to others.
        let /** @type {?} */ percentToDispatch = 0;
        // Get container pixel size
        let /** @type {?} */ containerSizePixel = this.getNbGutters() * this.gutterSize;
        if (this.direction === 'horizontal') {
            containerSizePixel = this.width ? this.width : this.elRef.nativeElement['offsetWidth'];
        }
        else {
            containerSizePixel = this.height ? this.height : this.elRef.nativeElement['offsetHeight'];
        }
        this.displayedAreas.forEach(area => {
            if (area.size * containerSizePixel < this.gutterSize) {
                percentToDispatch += area.size;
                area.size = 0;
            }
        });
        if (percentToDispatch > 0 && this.displayedAreas.length > 0) {
            const /** @type {?} */ nbAreasNotZero = this.displayedAreas.filter(a => a.size !== 0).length;
            if (nbAreasNotZero > 0) {
                const /** @type {?} */ percentToAdd = percentToDispatch / nbAreasNotZero;
                this.displayedAreas.filter(a => a.size !== 0).forEach(area => {
                    area.size += percentToAdd;
                });
            }
            else {
                this.displayedAreas[this.displayedAreas.length - 1].size = 1;
            }
        }
        this.refreshStyleSizes();
        this.cdRef.markForCheck();
    }
    /**
     * @return {?}
     */
    refreshStyleSizes() {
        const /** @type {?} */ sumGutterSize = this.getNbGutters() * this.gutterSize;
        this.displayedAreas.forEach(area => {
            area.comp.setStyleFlexbasis(`calc( ${area.size * 100}% - ${area.size * sumGutterSize}px )`, this.isDragging);
        });
    }
    /**
     * @param {?} startEvent
     * @param {?} gutterOrder
     * @param {?} gutterNum
     * @return {?}
     */
    startDragging(startEvent, gutterOrder, gutterNum) {
        startEvent.preventDefault();
        // Place code here to allow '(gutterClick)' event even if '[disabled]="true"'.
        this.currentGutterNum = gutterNum;
        this.draggingWithoutMove = true;
        this.ngZone.runOutsideAngular(() => {
            this.dragListeners.push(this.renderer.listen('document', 'mouseup', (e) => this.stopDragging()));
            this.dragListeners.push(this.renderer.listen('document', 'touchend', (e) => this.stopDragging()));
            this.dragListeners.push(this.renderer.listen('document', 'touchcancel', (e) => this.stopDragging()));
        });
        if (this.disabled) {
            return;
        }
        const /** @type {?} */ areaA = this.displayedAreas.find(a => a.order === gutterOrder - 1);
        const /** @type {?} */ areaB = this.displayedAreas.find(a => a.order === gutterOrder + 1);
        if (!areaA || !areaB) {
            return;
        }
        const /** @type {?} */ prop = (this.direction === 'horizontal') ? 'offsetWidth' : 'offsetHeight';
        this.dragStartValues.sizePixelContainer = this.elRef.nativeElement[prop];
        this.dragStartValues.sizePixelA = areaA.comp.getSizePixel(prop);
        this.dragStartValues.sizePixelB = areaB.comp.getSizePixel(prop);
        this.dragStartValues.sizePercentA = areaA.size;
        this.dragStartValues.sizePercentB = areaB.size;
        let /** @type {?} */ start;
        if (startEvent instanceof MouseEvent) {
            start = {
                x: startEvent.screenX,
                y: startEvent.screenY,
            };
        }
        else if (startEvent instanceof TouchEvent) {
            start = {
                x: startEvent.touches[0].screenX,
                y: startEvent.touches[0].screenY,
            };
        }
        else {
            return;
        }
        this.ngZone.runOutsideAngular(() => {
            this.dragListeners.push(this.renderer.listen('document', 'mousemove', (e) => this.dragEvent(e, start, areaA, areaB)));
            this.dragListeners.push(this.renderer.listen('document', 'touchmove', (e) => this.dragEvent(e, start, areaA, areaB)));
        });
        areaA.comp.lockEvents();
        areaB.comp.lockEvents();
        this.isDragging = true;
        this.notify('start');
    }
    /**
     * @param {?} event
     * @param {?} start
     * @param {?} areaA
     * @param {?} areaB
     * @return {?}
     */
    dragEvent(event, start, areaA, areaB) {
        if (!this.isDragging) {
            return;
        }
        let /** @type {?} */ end;
        if (event instanceof MouseEvent) {
            end = {
                x: event.screenX,
                y: event.screenY,
            };
        }
        else if (event instanceof TouchEvent) {
            end = {
                x: event.touches[0].screenX,
                y: event.touches[0].screenY,
            };
        }
        else {
            return;
        }
        this.draggingWithoutMove = false;
        this.drag(start, end, areaA, areaB);
    }
    /**
     * @param {?} start
     * @param {?} end
     * @param {?} areaA
     * @param {?} areaB
     * @return {?}
     */
    drag(start, end, areaA, areaB) {
        // ¤ AREAS SIZE PIXEL
        const /** @type {?} */ devicePixelRatio = window.devicePixelRatio || 1;
        let /** @type {?} */ offsetPixel = (this.direction === 'horizontal') ? (start.x - end.x) : (start.y - end.y);
        offsetPixel = offsetPixel / devicePixelRatio;
        if (this.dir === 'rtl') {
            offsetPixel = -offsetPixel;
        }
        let /** @type {?} */ newSizePixelA = this.dragStartValues.sizePixelA - offsetPixel;
        let /** @type {?} */ newSizePixelB = this.dragStartValues.sizePixelB + offsetPixel;
        if (newSizePixelA < this.gutterSize && newSizePixelB < this.gutterSize) {
            // WTF.. get out of here!
            return;
        }
        else if (newSizePixelA < this.gutterSize) {
            newSizePixelB += newSizePixelA;
            newSizePixelA = 0;
        }
        else if (newSizePixelB < this.gutterSize) {
            newSizePixelA += newSizePixelB;
            newSizePixelB = 0;
        }
        // ¤ AREAS SIZE PERCENT
        if (newSizePixelA === 0) {
            areaB.size += areaA.size;
            areaA.size = 0;
        }
        else if (newSizePixelB === 0) {
            areaA.size += areaB.size;
            areaB.size = 0;
        }
        else {
            // NEW_PERCENT = START_PERCENT / START_PIXEL * NEW_PIXEL;
            if (this.dragStartValues.sizePercentA === 0) {
                areaB.size = this.dragStartValues.sizePercentB / this.dragStartValues.sizePixelB * newSizePixelB;
                areaA.size = this.dragStartValues.sizePercentB - areaB.size;
            }
            else if (this.dragStartValues.sizePercentB === 0) {
                areaA.size = this.dragStartValues.sizePercentA / this.dragStartValues.sizePixelA * newSizePixelA;
                areaB.size = this.dragStartValues.sizePercentA - areaA.size;
            }
            else {
                areaA.size = this.dragStartValues.sizePercentA / this.dragStartValues.sizePixelA * newSizePixelA;
                areaB.size = (this.dragStartValues.sizePercentA + this.dragStartValues.sizePercentB) - areaA.size;
            }
        }
        this.refreshStyleSizes();
        this.notify('progress');
    }
    /**
     * @return {?}
     */
    stopDragging() {
        if (this.isDragging === false && this.draggingWithoutMove === false) {
            return;
        }
        this.displayedAreas.forEach(area => {
            area.comp.unlockEvents();
        });
        while (this.dragListeners.length > 0) {
            const /** @type {?} */ fct = this.dragListeners.pop();
            if (fct) {
                fct();
            }
        }
        if (this.draggingWithoutMove === true) {
            this.notify('click');
        }
        else {
            this.notify('end');
        }
        this.isDragging = false;
        this.draggingWithoutMove = false;
    }
    /**
     * @param {?} type
     * @return {?}
     */
    notify(type) {
        const /** @type {?} */ areasSize = this.displayedAreas.map(a => a.size * 100);
        switch (type) {
            case 'start':
                return this.dragStart.emit({ gutterNum: this.currentGutterNum, sizes: areasSize });
            case 'progress':
                return this.dragProgress.emit({ gutterNum: this.currentGutterNum, sizes: areasSize });
            case 'end':
                return this.dragEnd.emit({ gutterNum: this.currentGutterNum, sizes: areasSize });
            case 'click':
                return this.gutterClick.emit({ gutterNum: this.currentGutterNum, sizes: areasSize });
            case 'transitionEnd':
                return this.transitionEndInternal.next(areasSize);
        }
    }
    /**
     * @return {?}
     */
    ngOnDestroy() {
        this.stopDragging();
    }
}
SplitComponent.decorators = [
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Component"], args: [{
                selector: 'split',
                changeDetection: __WEBPACK_IMPORTED_MODULE_0__angular_core__["ChangeDetectionStrategy"].OnPush,
                styles: [`
        :host {
            display: flex;
            flex-wrap: nowrap;
            justify-content: flex-start;
            align-items: stretch;
            overflow: hidden;
            /* 
                Important to keep following rules even if overrided later by 'HostBinding' 
                because if [width] & [height] not provided, when build() is executed,
                'HostBinding' hasn't been applied yet so code:
                this.elRef.nativeElement["offsetHeight"] gives wrong value!  
             */
            width: 100%;
            height: 100%;   
        }

        split-gutter {
            flex-grow: 0;
            flex-shrink: 0;
            background-position: center center;
            background-repeat: no-repeat;
        }
    `],
                template: `
        <ng-content></ng-content>
        <ng-template ngFor let-area [ngForOf]="displayedAreas" let-index="index" let-last="last">
            <split-gutter *ngIf="last === false" 
                          [order]="index*2+1"
                          [direction]="direction"
                          [useTransition]="useTransition"
                          [size]="gutterSize"
                          [color]="gutterColor"
                          [imageH]="gutterImageH"
                          [imageV]="gutterImageV"
                          [disabled]="disabled"
                          (mousedown)="startDragging($event, index*2+1, index+1)"
                          (touchstart)="startDragging($event, index*2+1, index+1)"></split-gutter>
        </ng-template>`,
            },] },
];
/** @nocollapse */
SplitComponent.ctorParameters = () => [
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["NgZone"], },
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["ElementRef"], },
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["ChangeDetectorRef"], },
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Renderer2"], },
];
SplitComponent.propDecorators = {
    "direction": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "useTransition": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "disabled": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "width": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "height": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "gutterSize": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "gutterColor": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "gutterImageH": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "gutterImageV": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "dir": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "dragStart": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Output"] },],
    "dragProgress": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Output"] },],
    "dragEnd": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Output"] },],
    "gutterClick": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Output"] },],
    "transitionEnd": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Output"] },],
    "cssFlexdirection": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["HostBinding"], args: ['style.flex-direction',] },],
    "cssWidth": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["HostBinding"], args: ['style.width',] },],
    "cssHeight": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["HostBinding"], args: ['style.height',] },],
    "cssMinwidth": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["HostBinding"], args: ['style.min-width',] },],
    "cssMinheight": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["HostBinding"], args: ['style.min-height',] },],
};

/**
 * @fileoverview added by tsickle
 * @suppress {checkTypes} checked by tsc
 */
class SplitAreaDirective {
    /**
     * @param {?} ngZone
     * @param {?} elRef
     * @param {?} renderer
     * @param {?} split
     */
    constructor(ngZone, elRef, renderer, split) {
        this.ngZone = ngZone;
        this.elRef = elRef;
        this.renderer = renderer;
        this.split = split;
        this._order = null;
        this._size = null;
        this._minSize = 0;
        this._visible = true;
        this.lockListeners = [];
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set order(v) {
        v = Number(v);
        this._order = !isNaN(v) ? v : null;
        this.split.updateArea(this, true, false);
    }
    /**
     * @return {?}
     */
    get order() {
        return this._order;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set size(v) {
        v = Number(v);
        this._size = (!isNaN(v) && v >= 0 && v <= 100) ? (v / 100) : null;
        this.split.updateArea(this, false, true);
    }
    /**
     * @return {?}
     */
    get size() {
        return this._size;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set minSize(v) {
        v = Number(v);
        this._minSize = (!isNaN(v) && v > 0 && v < 100) ? v / 100 : 0;
        this.split.updateArea(this, false, true);
    }
    /**
     * @return {?}
     */
    get minSize() {
        return this._minSize;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set visible(v) {
        v = (typeof (v) === 'boolean') ? v : (v === 'false' ? false : true);
        this._visible = v;
        if (this.visible) {
            this.split.showArea(this);
        }
        else {
            this.split.hideArea(this);
        }
    }
    /**
     * @return {?}
     */
    get visible() {
        return this._visible;
    }
    /**
     * @return {?}
     */
    ngOnInit() {
        this.split.addArea(this);
        this.renderer.setStyle(this.elRef.nativeElement, 'flex-grow', '0');
        this.renderer.setStyle(this.elRef.nativeElement, 'flex-shrink', '0');
        this.ngZone.runOutsideAngular(() => {
            this.transitionListener = this.renderer.listen(this.elRef.nativeElement, 'transitionend', (e) => this.onTransitionEnd(e));
        });
    }
    /**
     * @param {?} prop
     * @return {?}
     */
    getSizePixel(prop) {
        return this.elRef.nativeElement[prop];
    }
    /**
     * @param {?} isVisible
     * @param {?} isDragging
     * @param {?} direction
     * @return {?}
     */
    setStyleVisibleAndDir(isVisible, isDragging, direction) {
        if (isVisible === false) {
            this.setStyleFlexbasis('0', isDragging);
            this.renderer.setStyle(this.elRef.nativeElement, 'overflow-x', 'hidden');
            this.renderer.setStyle(this.elRef.nativeElement, 'overflow-y', 'hidden');
            if (direction === 'vertical') {
                this.renderer.setStyle(this.elRef.nativeElement, 'max-width', '0');
            }
        }
        else {
            this.renderer.setStyle(this.elRef.nativeElement, 'overflow-x', 'hidden');
            this.renderer.setStyle(this.elRef.nativeElement, 'overflow-y', 'auto');
            this.renderer.removeStyle(this.elRef.nativeElement, 'max-width');
        }
        if (direction === 'horizontal') {
            this.renderer.setStyle(this.elRef.nativeElement, 'height', '100%');
            this.renderer.removeStyle(this.elRef.nativeElement, 'width');
        }
        else {
            this.renderer.setStyle(this.elRef.nativeElement, 'width', '100%');
            this.renderer.removeStyle(this.elRef.nativeElement, 'height');
        }
    }
    /**
     * @param {?} value
     * @return {?}
     */
    setStyleOrder(value) {
        this.renderer.setStyle(this.elRef.nativeElement, 'order', value);
    }
    /**
     * @param {?} value
     * @param {?} isDragging
     * @return {?}
     */
    setStyleFlexbasis(value, isDragging) {
        // If component not yet initialized or gutter being dragged, disable transition
        if (this.split.isViewInitialized === false || isDragging === true) {
            this.setStyleTransition(false);
        }
        else {
            this.setStyleTransition(this.split.useTransition);
        }
        this.renderer.setStyle(this.elRef.nativeElement, 'flex-basis', value);
    }
    /**
     * @param {?} useTransition
     * @return {?}
     */
    setStyleTransition(useTransition) {
        if (useTransition) {
            this.renderer.setStyle(this.elRef.nativeElement, 'transition', `flex-basis 0.3s`);
        }
        else {
            this.renderer.removeStyle(this.elRef.nativeElement, 'transition');
        }
    }
    /**
     * @param {?} event
     * @return {?}
     */
    onTransitionEnd(event) {
        // Limit only flex-basis transition to trigger the event
        if (event.propertyName === 'flex-basis') {
            this.split.notify('transitionEnd');
        }
    }
    /**
     * @return {?}
     */
    lockEvents() {
        this.ngZone.runOutsideAngular(() => {
            this.lockListeners.push(this.renderer.listen(this.elRef.nativeElement, 'selectstart', (e) => false));
            this.lockListeners.push(this.renderer.listen(this.elRef.nativeElement, 'dragstart', (e) => false));
        });
    }
    /**
     * @return {?}
     */
    unlockEvents() {
        while (this.lockListeners.length > 0) {
            const /** @type {?} */ fct = this.lockListeners.pop();
            if (fct) {
                fct();
            }
        }
    }
    /**
     * @return {?}
     */
    ngOnDestroy() {
        this.unlockEvents();
        if (this.transitionListener) {
            this.transitionListener();
        }
        this.split.removeArea(this);
    }
}
SplitAreaDirective.decorators = [
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Directive"], args: [{
                selector: 'split-area'
            },] },
];
/** @nocollapse */
SplitAreaDirective.ctorParameters = () => [
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["NgZone"], },
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["ElementRef"], },
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Renderer2"], },
    { type: SplitComponent, },
];
SplitAreaDirective.propDecorators = {
    "order": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "size": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "minSize": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "visible": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
};

/**
 * @fileoverview added by tsickle
 * @suppress {checkTypes} checked by tsc
 */
class SplitGutterDirective {
    /**
     * @param {?} elRef
     * @param {?} renderer
     */
    constructor(elRef, renderer) {
        this.elRef = elRef;
        this.renderer = renderer;
        this._disabled = false;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set order(v) {
        this.renderer.setStyle(this.elRef.nativeElement, 'order', v);
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set direction(v) {
        this._direction = v;
        this.refreshStyle();
    }
    /**
     * @return {?}
     */
    get direction() {
        return this._direction;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set useTransition(v) {
        if (v) {
            this.renderer.setStyle(this.elRef.nativeElement, 'transition', `flex-basis 0.3s`);
        }
        else {
            this.renderer.removeStyle(this.elRef.nativeElement, 'transition');
        }
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set size(v) {
        this._size = v;
        this.refreshStyle();
    }
    /**
     * @return {?}
     */
    get size() {
        return this._size;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set color(v) {
        this._color = v;
        this.refreshStyle();
    }
    /**
     * @return {?}
     */
    get color() {
        return this._color;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set imageH(v) {
        this._imageH = v;
        this.refreshStyle();
    }
    /**
     * @return {?}
     */
    get imageH() {
        return this._imageH;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set imageV(v) {
        this._imageV = v;
        this.refreshStyle();
    }
    /**
     * @return {?}
     */
    get imageV() {
        return this._imageV;
    }
    /**
     * @param {?} v
     * @return {?}
     */
    set disabled(v) {
        this._disabled = v;
        this.refreshStyle();
    }
    /**
     * @return {?}
     */
    get disabled() {
        return this._disabled;
    }
    /**
     * @return {?}
     */
    refreshStyle() {
        this.renderer.setStyle(this.elRef.nativeElement, 'flex-basis', `${this.size}px`);
        // fix safari bug about gutter height when direction is horizontal
        this.renderer.setStyle(this.elRef.nativeElement, 'height', (this.direction === 'vertical') ? `${this.size}px` : `100%`);
        this.renderer.setStyle(this.elRef.nativeElement, 'background-color', (this.color !== '') ? this.color : `#eeeeee`);
        const /** @type {?} */ state = (this.disabled === true) ? 'disabled' : this.direction;
        this.renderer.setStyle(this.elRef.nativeElement, 'background-image', this.getImage(state));
        this.renderer.setStyle(this.elRef.nativeElement, 'cursor', this.getCursor(state));
    }
    /**
     * @param {?} state
     * @return {?}
     */
    getCursor(state) {
        switch (state) {
            case 'horizontal':
                return 'col-resize';
            case 'vertical':
                return 'row-resize';
            case 'disabled':
                return 'default';
        }
    }
    /**
     * @param {?} state
     * @return {?}
     */
    getImage(state) {
        switch (state) {
            case 'horizontal':
                return (this.imageH !== '') ? this.imageH : defaultImageH;
            case 'vertical':
                return (this.imageV !== '') ? this.imageV : defaultImageV;
            case 'disabled':
                return '';
        }
    }
}
SplitGutterDirective.decorators = [
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Directive"], args: [{
                selector: 'split-gutter'
            },] },
];
/** @nocollapse */
SplitGutterDirective.ctorParameters = () => [
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["ElementRef"], },
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Renderer2"], },
];
SplitGutterDirective.propDecorators = {
    "order": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "direction": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "useTransition": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "size": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "color": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "imageH": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "imageV": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
    "disabled": [{ type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["Input"] },],
};
const defaultImageH = 'url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAeCAYAAADkftS9AAAAIklEQVQoU2M4c+bMfxAGAgYYmwGrIIiDjrELjpo5aiZeMwF+yNnOs5KSvgAAAABJRU5ErkJggg==")';
const defaultImageV = 'url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAB4AAAAFCAMAAABl/6zIAAAABlBMVEUAAADMzMzIT8AyAAAAAXRSTlMAQObYZgAAABRJREFUeAFjYGRkwIMJSeMHlBkOABP7AEGzSuPKAAAAAElFTkSuQmCC")';

/**
 * @fileoverview added by tsickle
 * @suppress {checkTypes} checked by tsc
 */
class AngularSplitModule {
    /**
     * @return {?}
     */
    static forRoot() {
        return {
            ngModule: AngularSplitModule,
            providers: []
        };
    }
    /**
     * @return {?}
     */
    static forChild() {
        return {
            ngModule: AngularSplitModule,
            providers: []
        };
    }
}
AngularSplitModule.decorators = [
    { type: __WEBPACK_IMPORTED_MODULE_0__angular_core__["NgModule"], args: [{
                imports: [
                    __WEBPACK_IMPORTED_MODULE_1__angular_common__["CommonModule"]
                ],
                declarations: [
                    SplitComponent,
                    SplitAreaDirective,
                    SplitGutterDirective,
                ],
                exports: [
                    SplitComponent,
                    SplitAreaDirective,
                ]
            },] },
];
/** @nocollapse */
AngularSplitModule.ctorParameters = () => [];

/**
 * @fileoverview added by tsickle
 * @suppress {checkTypes} checked by tsc
 */
// Public classes.

/**
 * @fileoverview added by tsickle
 * @suppress {checkTypes} checked by tsc
 */
/**
 * Entry point for all public APIs of the package.
 */

/**
 * @fileoverview added by tsickle
 * @suppress {checkTypes} checked by tsc
 */
/**
 * Generated bundle index. Do not edit.
 */


//# sourceMappingURL=angular-split.js.map


/***/ }),

/***/ "./node_modules/p-queue/index.js":
/***/ (function(module, exports, __webpack_require__) {

"use strict";


// Port of lower_bound from http://en.cppreference.com/w/cpp/algorithm/lower_bound
// Used to compute insertion index to keep queue sorted after insertion
function lowerBound(array, value, comp) {
	let first = 0;
	let count = array.length;

	while (count > 0) {
		const step = (count / 2) | 0;
		let it = first + step;

		if (comp(array[it], value) <= 0) {
			first = ++it;
			count -= step + 1;
		} else {
			count = step;
		}
	}

	return first;
}

class PriorityQueue {
	constructor() {
		this._queue = [];
	}

	enqueue(run, opts) {
		opts = Object.assign({
			priority: 0
		}, opts);

		const element = {priority: opts.priority, run};

		if (this.size && this._queue[this.size - 1].priority >= opts.priority) {
			this._queue.push(element);
			return;
		}

		const index = lowerBound(this._queue, element, (a, b) => b.priority - a.priority);
		this._queue.splice(index, 0, element);
	}

	dequeue() {
		return this._queue.shift().run;
	}

	get size() {
		return this._queue.length;
	}
}

class PQueue {
	constructor(opts) {
		opts = Object.assign({
			concurrency: Infinity,
			autoStart: true,
			queueClass: PriorityQueue
		}, opts);

		if (!(typeof opts.concurrency === 'number' && opts.concurrency >= 1)) {
			throw new TypeError(`Expected \`concurrency\` to be a number from 1 and up, got \`${opts.concurrency}\` (${typeof opts.concurrency})`);
		}

		this.queue = new opts.queueClass(); // eslint-disable-line new-cap
		this._queueClass = opts.queueClass;
		this._pendingCount = 0;
		this._concurrency = opts.concurrency;
		this._isPaused = opts.autoStart === false;
		this._resolveEmpty = () => {};
		this._resolveIdle = () => {};
	}

	_next() {
		this._pendingCount--;

		if (this.queue.size > 0) {
			if (!this._isPaused) {
				this.queue.dequeue()();
			}
		} else {
			this._resolveEmpty();
			this._resolveEmpty = () => {};

			if (this._pendingCount === 0) {
				this._resolveIdle();
				this._resolveIdle = () => {};
			}
		}
	}

	add(fn, opts) {
		return new Promise((resolve, reject) => {
			const run = () => {
				this._pendingCount++;

				try {
					Promise.resolve(fn()).then(
						val => {
							resolve(val);
							this._next();
						},
						err => {
							reject(err);
							this._next();
						}
					);
				} catch (err) {
					reject(err);
					this._next();
				}
			};

			if (!this._isPaused && this._pendingCount < this._concurrency) {
				run();
			} else {
				this.queue.enqueue(run, opts);
			}
		});
	}

	addAll(fns, opts) {
		return Promise.all(fns.map(fn => this.add(fn, opts)));
	}

	start() {
		if (!this._isPaused) {
			return;
		}

		this._isPaused = false;
		while (this.queue.size > 0 && this._pendingCount < this._concurrency) {
			this.queue.dequeue()();
		}
	}

	pause() {
		this._isPaused = true;
	}

	clear() {
		this.queue = new this._queueClass(); // eslint-disable-line new-cap
	}

	onEmpty() {
		// Instantly resolve if the queue is empty
		if (this.queue.size === 0) {
			return Promise.resolve();
		}

		return new Promise(resolve => {
			const existingResolve = this._resolveEmpty;
			this._resolveEmpty = () => {
				existingResolve();
				resolve();
			};
		});
	}

	onIdle() {
		// Instantly resolve if none pending
		if (this._pendingCount === 0) {
			return Promise.resolve();
		}

		return new Promise(resolve => {
			const existingResolve = this._resolveIdle;
			this._resolveIdle = () => {
				existingResolve();
				resolve();
			};
		});
	}

	get size() {
		return this.queue.size;
	}

	get pending() {
		return this._pendingCount;
	}

	get isPaused() {
		return this._isPaused;
	}
}

module.exports = PQueue;


/***/ }),

/***/ "./node_modules/primeng/api.js":
/***/ (function(module, exports, __webpack_require__) {

"use strict";
/* Shorthand */

function __export(m) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });
__export(__webpack_require__("./node_modules/primeng/components/common/api.js"));

/***/ }),

/***/ "./node_modules/primeng/button.js":
/***/ (function(module, exports, __webpack_require__) {

"use strict";
/* Shorthand */

function __export(m) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });
__export(__webpack_require__("./node_modules/primeng/components/button/button.js"));

/***/ }),

/***/ "./node_modules/primeng/dragdrop.js":
/***/ (function(module, exports, __webpack_require__) {

"use strict";
/* Shorthand */

function __export(m) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });
__export(__webpack_require__("./node_modules/primeng/components/dragdrop/dragdrop.js"));

/***/ }),

/***/ "./node_modules/primeng/inputtext.js":
/***/ (function(module, exports, __webpack_require__) {

"use strict";
/* Shorthand */

function __export(m) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });
__export(__webpack_require__("./node_modules/primeng/components/inputtext/inputtext.js"));

/***/ }),

/***/ "./node_modules/primeng/radiobutton.js":
/***/ (function(module, exports, __webpack_require__) {

"use strict";
/* Shorthand */

function __export(m) {
    for (var p in m) if (!exports.hasOwnProperty(p)) exports[p] = m[p];
}
Object.defineProperty(exports, "__esModule", { value: true });
__export(__webpack_require__("./node_modules/primeng/components/radiobutton/radiobutton.js"));

/***/ }),

/***/ "./node_modules/rxjs/_esm2015/Scheduler.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/**
 * An execution context and a data structure to order tasks and schedule their
 * execution. Provides a notion of (potentially virtual) time, through the
 * `now()` getter method.
 *
 * Each unit of work in a Scheduler is called an {@link Action}.
 *
 * ```ts
 * class Scheduler {
 *   now(): number;
 *   schedule(work, delay?, state?): Subscription;
 * }
 * ```
 *
 * @class Scheduler
 */
class Scheduler {
    constructor(SchedulerAction, now = Scheduler.now) {
        this.SchedulerAction = SchedulerAction;
        this.now = now;
    }
    /**
     * Schedules a function, `work`, for execution. May happen at some point in
     * the future, according to the `delay` parameter, if specified. May be passed
     * some context object, `state`, which will be passed to the `work` function.
     *
     * The given arguments will be processed an stored as an Action object in a
     * queue of actions.
     *
     * @param {function(state: ?T): ?Subscription} work A function representing a
     * task, or some unit of work to be executed by the Scheduler.
     * @param {number} [delay] Time to wait before executing the work, where the
     * time unit is implicit and defined by the Scheduler itself.
     * @param {T} [state] Some contextual data that the `work` function uses when
     * called by the Scheduler.
     * @return {Subscription} A subscription in order to be able to unsubscribe
     * the scheduled work.
     */
    schedule(work, delay = 0, state) {
        return new this.SchedulerAction(this, work).schedule(state, delay);
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = Scheduler;

Scheduler.now = Date.now ? Date.now : () => +new Date();
//# sourceMappingURL=Scheduler.js.map

/***/ }),

/***/ "./node_modules/rxjs/_esm2015/add/operator/debounceTime.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__Observable__ = __webpack_require__("./node_modules/rxjs/_esm2015/Observable.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__operator_debounceTime__ = __webpack_require__("./node_modules/rxjs/_esm2015/operator/debounceTime.js");


__WEBPACK_IMPORTED_MODULE_0__Observable__["Observable"].prototype.debounceTime = __WEBPACK_IMPORTED_MODULE_1__operator_debounceTime__["a" /* debounceTime */];
//# sourceMappingURL=debounceTime.js.map

/***/ }),

/***/ "./node_modules/rxjs/_esm2015/operator/debounceTime.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (immutable) */ __webpack_exports__["a"] = debounceTime;
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__scheduler_async__ = __webpack_require__("./node_modules/rxjs/_esm2015/scheduler/async.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__operators_debounceTime__ = __webpack_require__("./node_modules/rxjs/_esm2015/operators/debounceTime.js");


/**
 * Emits a value from the source Observable only after a particular time span
 * has passed without another source emission.
 *
 * <span class="informal">It's like {@link delay}, but passes only the most
 * recent value from each burst of emissions.</span>
 *
 * <img src="./img/debounceTime.png" width="100%">
 *
 * `debounceTime` delays values emitted by the source Observable, but drops
 * previous pending delayed emissions if a new value arrives on the source
 * Observable. This operator keeps track of the most recent value from the
 * source Observable, and emits that only when `dueTime` enough time has passed
 * without any other value appearing on the source Observable. If a new value
 * appears before `dueTime` silence occurs, the previous value will be dropped
 * and will not be emitted on the output Observable.
 *
 * This is a rate-limiting operator, because it is impossible for more than one
 * value to be emitted in any time window of duration `dueTime`, but it is also
 * a delay-like operator since output emissions do not occur at the same time as
 * they did on the source Observable. Optionally takes a {@link IScheduler} for
 * managing timers.
 *
 * @example <caption>Emit the most recent click after a burst of clicks</caption>
 * var clicks = Rx.Observable.fromEvent(document, 'click');
 * var result = clicks.debounceTime(1000);
 * result.subscribe(x => console.log(x));
 *
 * @see {@link auditTime}
 * @see {@link debounce}
 * @see {@link delay}
 * @see {@link sampleTime}
 * @see {@link throttleTime}
 *
 * @param {number} dueTime The timeout duration in milliseconds (or the time
 * unit determined internally by the optional `scheduler`) for the window of
 * time required to wait for emission silence before emitting the most recent
 * source value.
 * @param {Scheduler} [scheduler=async] The {@link IScheduler} to use for
 * managing the timers that handle the timeout for each value.
 * @return {Observable} An Observable that delays the emissions of the source
 * Observable by the specified `dueTime`, and may drop some values if they occur
 * too frequently.
 * @method debounceTime
 * @owner Observable
 */
function debounceTime(dueTime, scheduler = __WEBPACK_IMPORTED_MODULE_0__scheduler_async__["a" /* async */]) {
    return Object(__WEBPACK_IMPORTED_MODULE_1__operators_debounceTime__["a" /* debounceTime */])(dueTime, scheduler)(this);
}
//# sourceMappingURL=debounceTime.js.map

/***/ }),

/***/ "./node_modules/rxjs/_esm2015/operators/debounceTime.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony export (immutable) */ __webpack_exports__["a"] = debounceTime;
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__Subscriber__ = __webpack_require__("./node_modules/rxjs/_esm2015/Subscriber.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__scheduler_async__ = __webpack_require__("./node_modules/rxjs/_esm2015/scheduler/async.js");


/**
 * Emits a value from the source Observable only after a particular time span
 * has passed without another source emission.
 *
 * <span class="informal">It's like {@link delay}, but passes only the most
 * recent value from each burst of emissions.</span>
 *
 * <img src="./img/debounceTime.png" width="100%">
 *
 * `debounceTime` delays values emitted by the source Observable, but drops
 * previous pending delayed emissions if a new value arrives on the source
 * Observable. This operator keeps track of the most recent value from the
 * source Observable, and emits that only when `dueTime` enough time has passed
 * without any other value appearing on the source Observable. If a new value
 * appears before `dueTime` silence occurs, the previous value will be dropped
 * and will not be emitted on the output Observable.
 *
 * This is a rate-limiting operator, because it is impossible for more than one
 * value to be emitted in any time window of duration `dueTime`, but it is also
 * a delay-like operator since output emissions do not occur at the same time as
 * they did on the source Observable. Optionally takes a {@link IScheduler} for
 * managing timers.
 *
 * @example <caption>Emit the most recent click after a burst of clicks</caption>
 * var clicks = Rx.Observable.fromEvent(document, 'click');
 * var result = clicks.debounceTime(1000);
 * result.subscribe(x => console.log(x));
 *
 * @see {@link auditTime}
 * @see {@link debounce}
 * @see {@link delay}
 * @see {@link sampleTime}
 * @see {@link throttleTime}
 *
 * @param {number} dueTime The timeout duration in milliseconds (or the time
 * unit determined internally by the optional `scheduler`) for the window of
 * time required to wait for emission silence before emitting the most recent
 * source value.
 * @param {Scheduler} [scheduler=async] The {@link IScheduler} to use for
 * managing the timers that handle the timeout for each value.
 * @return {Observable} An Observable that delays the emissions of the source
 * Observable by the specified `dueTime`, and may drop some values if they occur
 * too frequently.
 * @method debounceTime
 * @owner Observable
 */
function debounceTime(dueTime, scheduler = __WEBPACK_IMPORTED_MODULE_1__scheduler_async__["a" /* async */]) {
    return (source) => source.lift(new DebounceTimeOperator(dueTime, scheduler));
}
class DebounceTimeOperator {
    constructor(dueTime, scheduler) {
        this.dueTime = dueTime;
        this.scheduler = scheduler;
    }
    call(subscriber, source) {
        return source.subscribe(new DebounceTimeSubscriber(subscriber, this.dueTime, this.scheduler));
    }
}
/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
class DebounceTimeSubscriber extends __WEBPACK_IMPORTED_MODULE_0__Subscriber__["a" /* Subscriber */] {
    constructor(destination, dueTime, scheduler) {
        super(destination);
        this.dueTime = dueTime;
        this.scheduler = scheduler;
        this.debouncedSubscription = null;
        this.lastValue = null;
        this.hasValue = false;
    }
    _next(value) {
        this.clearDebounce();
        this.lastValue = value;
        this.hasValue = true;
        this.add(this.debouncedSubscription = this.scheduler.schedule(dispatchNext, this.dueTime, this));
    }
    _complete() {
        this.debouncedNext();
        this.destination.complete();
    }
    debouncedNext() {
        this.clearDebounce();
        if (this.hasValue) {
            this.destination.next(this.lastValue);
            this.lastValue = null;
            this.hasValue = false;
        }
    }
    clearDebounce() {
        const debouncedSubscription = this.debouncedSubscription;
        if (debouncedSubscription !== null) {
            this.remove(debouncedSubscription);
            debouncedSubscription.unsubscribe();
            this.debouncedSubscription = null;
        }
    }
}
function dispatchNext(subscriber) {
    subscriber.debouncedNext();
}
//# sourceMappingURL=debounceTime.js.map

/***/ }),

/***/ "./node_modules/rxjs/_esm2015/scheduler/Action.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__Subscription__ = __webpack_require__("./node_modules/rxjs/_esm2015/Subscription.js");

/**
 * A unit of work to be executed in a {@link Scheduler}. An action is typically
 * created from within a Scheduler and an RxJS user does not need to concern
 * themselves about creating and manipulating an Action.
 *
 * ```ts
 * class Action<T> extends Subscription {
 *   new (scheduler: Scheduler, work: (state?: T) => void);
 *   schedule(state?: T, delay: number = 0): Subscription;
 * }
 * ```
 *
 * @class Action<T>
 */
class Action extends __WEBPACK_IMPORTED_MODULE_0__Subscription__["a" /* Subscription */] {
    constructor(scheduler, work) {
        super();
    }
    /**
     * Schedules this action on its parent Scheduler for execution. May be passed
     * some context object, `state`. May happen at some point in the future,
     * according to the `delay` parameter, if specified.
     * @param {T} [state] Some contextual data that the `work` function uses when
     * called by the Scheduler.
     * @param {number} [delay] Time to wait before executing the work, where the
     * time unit is implicit and defined by the Scheduler.
     * @return {void}
     */
    schedule(state, delay = 0) {
        return this;
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = Action;

//# sourceMappingURL=Action.js.map

/***/ }),

/***/ "./node_modules/rxjs/_esm2015/scheduler/AsyncAction.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__util_root__ = __webpack_require__("./node_modules/rxjs/_esm2015/util/root.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__Action__ = __webpack_require__("./node_modules/rxjs/_esm2015/scheduler/Action.js");


/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
class AsyncAction extends __WEBPACK_IMPORTED_MODULE_1__Action__["a" /* Action */] {
    constructor(scheduler, work) {
        super(scheduler, work);
        this.scheduler = scheduler;
        this.work = work;
        this.pending = false;
    }
    schedule(state, delay = 0) {
        if (this.closed) {
            return this;
        }
        // Always replace the current state with the new state.
        this.state = state;
        // Set the pending flag indicating that this action has been scheduled, or
        // has recursively rescheduled itself.
        this.pending = true;
        const id = this.id;
        const scheduler = this.scheduler;
        //
        // Important implementation note:
        //
        // Actions only execute once by default, unless rescheduled from within the
        // scheduled callback. This allows us to implement single and repeat
        // actions via the same code path, without adding API surface area, as well
        // as mimic traditional recursion but across asynchronous boundaries.
        //
        // However, JS runtimes and timers distinguish between intervals achieved by
        // serial `setTimeout` calls vs. a single `setInterval` call. An interval of
        // serial `setTimeout` calls can be individually delayed, which delays
        // scheduling the next `setTimeout`, and so on. `setInterval` attempts to
        // guarantee the interval callback will be invoked more precisely to the
        // interval period, regardless of load.
        //
        // Therefore, we use `setInterval` to schedule single and repeat actions.
        // If the action reschedules itself with the same delay, the interval is not
        // canceled. If the action doesn't reschedule, or reschedules with a
        // different delay, the interval will be canceled after scheduled callback
        // execution.
        //
        if (id != null) {
            this.id = this.recycleAsyncId(scheduler, id, delay);
        }
        this.delay = delay;
        // If this action has already an async Id, don't request a new one.
        this.id = this.id || this.requestAsyncId(scheduler, this.id, delay);
        return this;
    }
    requestAsyncId(scheduler, id, delay = 0) {
        return __WEBPACK_IMPORTED_MODULE_0__util_root__["a" /* root */].setInterval(scheduler.flush.bind(scheduler, this), delay);
    }
    recycleAsyncId(scheduler, id, delay = 0) {
        // If this action is rescheduled with the same delay time, don't clear the interval id.
        if (delay !== null && this.delay === delay && this.pending === false) {
            return id;
        }
        // Otherwise, if the action's delay time is different from the current delay,
        // or the action has been rescheduled before it's executed, clear the interval id
        return __WEBPACK_IMPORTED_MODULE_0__util_root__["a" /* root */].clearInterval(id) && undefined || undefined;
    }
    /**
     * Immediately executes this action and the `work` it contains.
     * @return {any}
     */
    execute(state, delay) {
        if (this.closed) {
            return new Error('executing a cancelled action');
        }
        this.pending = false;
        const error = this._execute(state, delay);
        if (error) {
            return error;
        }
        else if (this.pending === false && this.id != null) {
            // Dequeue if the action didn't reschedule itself. Don't call
            // unsubscribe(), because the action could reschedule later.
            // For example:
            // ```
            // scheduler.schedule(function doWork(counter) {
            //   /* ... I'm a busy worker bee ... */
            //   var originalAction = this;
            //   /* wait 100ms before rescheduling the action */
            //   setTimeout(function () {
            //     originalAction.schedule(counter + 1);
            //   }, 100);
            // }, 1000);
            // ```
            this.id = this.recycleAsyncId(this.scheduler, this.id, null);
        }
    }
    _execute(state, delay) {
        let errored = false;
        let errorValue = undefined;
        try {
            this.work(state);
        }
        catch (e) {
            errored = true;
            errorValue = !!e && e || new Error(e);
        }
        if (errored) {
            this.unsubscribe();
            return errorValue;
        }
    }
    /** @deprecated internal use only */ _unsubscribe() {
        const id = this.id;
        const scheduler = this.scheduler;
        const actions = scheduler.actions;
        const index = actions.indexOf(this);
        this.work = null;
        this.state = null;
        this.pending = false;
        this.scheduler = null;
        if (index !== -1) {
            actions.splice(index, 1);
        }
        if (id != null) {
            this.id = this.recycleAsyncId(scheduler, id, null);
        }
        this.delay = null;
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = AsyncAction;

//# sourceMappingURL=AsyncAction.js.map

/***/ }),

/***/ "./node_modules/rxjs/_esm2015/scheduler/AsyncScheduler.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__Scheduler__ = __webpack_require__("./node_modules/rxjs/_esm2015/Scheduler.js");

class AsyncScheduler extends __WEBPACK_IMPORTED_MODULE_0__Scheduler__["a" /* Scheduler */] {
    constructor() {
        super(...arguments);
        this.actions = [];
        /**
         * A flag to indicate whether the Scheduler is currently executing a batch of
         * queued actions.
         * @type {boolean}
         */
        this.active = false;
        /**
         * An internal ID used to track the latest asynchronous task such as those
         * coming from `setTimeout`, `setInterval`, `requestAnimationFrame`, and
         * others.
         * @type {any}
         */
        this.scheduled = undefined;
    }
    flush(action) {
        const { actions } = this;
        if (this.active) {
            actions.push(action);
            return;
        }
        let error;
        this.active = true;
        do {
            if (error = action.execute(action.state, action.delay)) {
                break;
            }
        } while (action = actions.shift()); // exhaust the scheduler queue
        this.active = false;
        if (error) {
            while (action = actions.shift()) {
                action.unsubscribe();
            }
            throw error;
        }
    }
}
/* harmony export (immutable) */ __webpack_exports__["a"] = AsyncScheduler;

//# sourceMappingURL=AsyncScheduler.js.map

/***/ }),

/***/ "./node_modules/rxjs/_esm2015/scheduler/async.js":
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__AsyncAction__ = __webpack_require__("./node_modules/rxjs/_esm2015/scheduler/AsyncAction.js");
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_1__AsyncScheduler__ = __webpack_require__("./node_modules/rxjs/_esm2015/scheduler/AsyncScheduler.js");


/**
 *
 * Async Scheduler
 *
 * <span class="informal">Schedule task as if you used setTimeout(task, duration)</span>
 *
 * `async` scheduler schedules tasks asynchronously, by putting them on the JavaScript
 * event loop queue. It is best used to delay tasks in time or to schedule tasks repeating
 * in intervals.
 *
 * If you just want to "defer" task, that is to perform it right after currently
 * executing synchronous code ends (commonly achieved by `setTimeout(deferredTask, 0)`),
 * better choice will be the {@link asap} scheduler.
 *
 * @example <caption>Use async scheduler to delay task</caption>
 * const task = () => console.log('it works!');
 *
 * Rx.Scheduler.async.schedule(task, 2000);
 *
 * // After 2 seconds logs:
 * // "it works!"
 *
 *
 * @example <caption>Use async scheduler to repeat task in intervals</caption>
 * function task(state) {
 *   console.log(state);
 *   this.schedule(state + 1, 1000); // `this` references currently executing Action,
 *                                   // which we reschedule with new state and delay
 * }
 *
 * Rx.Scheduler.async.schedule(task, 3000, 0);
 *
 * // Logs:
 * // 0 after 3s
 * // 1 after 4s
 * // 2 after 5s
 * // 3 after 6s
 *
 * @static true
 * @name async
 * @owner Scheduler
 */
const async = new __WEBPACK_IMPORTED_MODULE_1__AsyncScheduler__["a" /* AsyncScheduler */](__WEBPACK_IMPORTED_MODULE_0__AsyncAction__["a" /* AsyncAction */]);
/* harmony export (immutable) */ __webpack_exports__["a"] = async;

//# sourceMappingURL=async.js.map

/***/ }),

/***/ "./src/app/ontology.conversation.history/conversation-history.component.html":
/***/ (function(module, exports) {

module.exports = "<div *ngIf=\"conversationHistory?.length\">\n  <div *ngFor=\"let message of conversationHistory; trackBy: message?.dialogState?.sn\" class=\"oty-bubble-speech\">\n    <div class=\"flex-container\" *ngIf=\"message.content.botResponseCaption?.length\">\n      <div class=\"oty-bubble-user\" *ngFor=\"let item of message.content.botResponseCaption\">\n        <span [title]=\"message.dialogState.timestamp | date:'short'\">{{item}}</span>\n      </div>\n    </div>\n\n    <div class=\"flex-container\" *ngIf=\"message.content.botGreetingCaption?.length\">\n      <div class=\"oty-robot-icon\">\n        <img src=\"../../assets/images/icon_antenna_bot.png\">\n      </div>\n      <div class=\"oty-bubble-container\">\n        <div class=\"oty-bubble-robot\" *ngFor=\"let item of message.content.botGreetingCaption\">\n          <span [title]=\"message.dialogState.timestamp | date:'short'\" >{{item}}</span>\n        </div>\n      </div>\n    </div>\n\n  </div>\n  <section class=\"oty-wrapper-options\" *ngIf=\"multiselectOptions?.length\">\n    <label *ngFor=\"let option of multiselectOptions\">\n      <input type=\"checkbox\" class=\"oty-options\" [(ngModel)]=\"option.checked\"/>{{option.text}}\n    </label>\n    <button (click)=\"sendMultiOptions()\">Done</button>\n  </section>\n\n  <section class=\"oty-wrapper-options\" *ngIf=\"buttonActions?.length\">\n    <button *ngFor=\"let option of buttonActions\" [ngClass]=\"{'disable': option.disable}\"\n            (click)=\"option.disable ? '' : sendButtonOption(option.value)\" [disabled]=\"option.disable\" class=\"oty-options\">{{option.label}}\n    </button>\n  </section>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.conversation.history/conversation-history.component.scss":
/***/ (function(module, exports) {

module.exports = ":host {\n  overflow: auto;\n  padding: 10px 0 40px; }\n  :host .oty-bubble-speech .flex-container {\n    display: -webkit-box;\n    display: -ms-flexbox;\n    display: flex; }\n  :host .oty-bubble-speech .flex-container .oty-bubble-container {\n      display: -webkit-box;\n      display: -ms-flexbox;\n      display: flex;\n      -webkit-box-orient: vertical;\n      -webkit-box-direction: normal;\n          -ms-flex-direction: column;\n              flex-direction: column;\n      width: 100%; }\n  :host .oty-bubble-speech .flex-container .oty-bubble-user,\n    :host .oty-bubble-speech .flex-container .oty-bubble-robot {\n      border-radius: 0 17.5px 17.5px 17.5px;\n      padding: 10px;\n      overflow-wrap: break-word;\n      margin-bottom: 4px;\n      min-width: 60px; }\n  :host .oty-bubble-speech .flex-container .oty-bubble-user span,\n      :host .oty-bubble-speech .flex-container .oty-bubble-robot span {\n        display: block; }\n  :host .oty-bubble-speech .flex-container .oty-bubble-user {\n      margin: 16px 20px 16px auto;\n      border-radius: 17.5px 0 17.5px 17.5px;\n      background-color: #E1F5FE;\n      max-width: 60%; }\n  :host .oty-bubble-speech .flex-container .oty-robot-icon {\n      margin-left: 24px; }\n  :host .oty-bubble-speech .flex-container .oty-robot-icon img {\n        width: 32px;\n        height: 32px;\n        border-radius: 50%;\n        border: 1px #DEDEDE solid; }\n  :host .oty-bubble-speech .flex-container .oty-bubble-robot {\n      width: 70%;\n      max-width: 80%;\n      margin-left: 9px;\n      background-color: #F4F4F4;\n      border: 1px #F4F4F4 solid; }\n  :host .oty-wrapper-options {\n    text-align: center;\n    padding-right: 15px; }\n  :host .oty-wrapper-options > button {\n      margin: 10px;\n      padding: 10px;\n      background: #FFFFFF;\n      width: 120px;\n      border-radius: 5%;\n      border: 1px #0081cb solid;\n      color: #0081cb; }\n  :host .oty-wrapper-options > button:hover {\n        color: #0051AC;\n        border: solid 1px #0051AC; }\n  :host .oty-wrapper-options > button:active {\n        background-color: #E1F5FE;\n        border: none;\n        border-radius: 17.5px 0 17.5px 17.5px;\n        color: #10203F; }\n  :host .oty-wrapper-options > button:focus {\n        outline: none; }\n  :host .oty-wrapper-options .disable {\n      color: rgba(0, 129, 203, 0.47) !important;\n      border-color: rgba(0, 129, 203, 0.47) !important; }\n  :host .oty-wrapper-options label {\n      display: block;\n      font-weight: normal;\n      text-align: left;\n      padding: 2px 0 2px 20px; }\n  :host .oty-wrapper-options label input {\n        margin-right: 10px; }\n"

/***/ }),

/***/ "./src/app/ontology.conversation.history/conversation-history.component.ts":
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
let ConversationHistoryComponent = class ConversationHistoryComponent {
    constructor(myElement, conversationService) {
        this.myElement = myElement;
        this.conversationService = conversationService;
        this.subscriptions = [];
        this.onSendOption = new core_1.EventEmitter();
        this.onSendMultiOptions = new core_1.EventEmitter();
    }
    ngOnInit() {
        this.subscriptions.push(this.conversationService.scrollConversationHistoryList.subscribe(() => {
            this.scrollToBottom();
        }));
    }
    ngOnDestroy() {
        this.subscriptions.forEach((subscription) => subscription.unsubscribe());
    }
    scrollToBottom() {
        setTimeout(() => {
            this.myElement.nativeElement.scrollTop = this.myElement.nativeElement.scrollHeight;
        }, 0);
    }
    sendButtonOption(option) {
        this.onSendOption.emit(option);
    }
    sendMultiOptions() {
        let selectedItem = [];
        this.multiselectOptions.forEach((item) => {
            if (item.checked) {
                selectedItem.push(item.value);
            }
        });
        this.onSendMultiOptions.emit(selectedItem);
    }
};
__decorate([
    core_1.Input(),
    __metadata("design:type", Array)
], ConversationHistoryComponent.prototype, "conversationHistory", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Array)
], ConversationHistoryComponent.prototype, "buttonActions", void 0);
__decorate([
    core_1.Input(),
    __metadata("design:type", Array)
], ConversationHistoryComponent.prototype, "multiselectOptions", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], ConversationHistoryComponent.prototype, "onSendOption", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], ConversationHistoryComponent.prototype, "onSendMultiOptions", void 0);
ConversationHistoryComponent = __decorate([
    core_1.Component({
        selector: 'oty-conversation-history',
        template: __webpack_require__("./src/app/ontology.conversation.history/conversation-history.component.html"),
        styles: [__webpack_require__("./src/app/ontology.conversation.history/conversation-history.component.scss"), __webpack_require__("./src/app/ontology.conversation/conversation.component.main.scss")]
    }),
    __metadata("design:paramtypes", [core_1.ElementRef,
        conversation_service_1.ConversationService])
], ConversationHistoryComponent);
exports.ConversationHistoryComponent = ConversationHistoryComponent;


/***/ }),

/***/ "./src/app/ontology.conversation.history/conversation-history.module.ts":
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
const conversation_history_component_1 = __webpack_require__("./src/app/ontology.conversation.history/conversation-history.component.ts");
const shared_module_1 = __webpack_require__("./src/app/common/shared.module.ts");
let ConversationHistoryModule = class ConversationHistoryModule {
};
ConversationHistoryModule = __decorate([
    core_1.NgModule({
        imports: [common_1.CommonModule, forms_1.FormsModule, shared_module_1.SharedModule],
        declarations: [conversation_history_component_1.ConversationHistoryComponent],
        exports: [conversation_history_component_1.ConversationHistoryComponent]
    })
], ConversationHistoryModule);
exports.ConversationHistoryModule = ConversationHistoryModule;


/***/ }),

/***/ "./src/app/ontology.conversation/conversation.component.html":
/***/ (function(module, exports) {

module.exports = "<div id=\"oty-conversation\" [class.oty-with-input]=\"inputSchema\">\n  <oty-conversation-history [conversationHistory]=\"getConversationHistoryList()\"\n                            [buttonActions]=\"buttonActions\"\n                            [multiselectOptions]=\"multiselectOptions\"\n                            (onSendOption)=\"sendButtonOption($event)\"\n                            (onSendMultiOptions)=\"sendMultiOptions($event)\"></oty-conversation-history>\n\n  <div class=\"oty-conversation-form\" [class.is-list]=\"isListAction\">\n    <form [formGroup]=\"elemForm\" (ngSubmit)=\"sendTyping(elemForm, isListAction)\">\n      <input type=\"text\" pInputText autofocus placeholder=\"{{inputPlaceholder}}\"\n             formControlName=\"userInput\" required/>\n      <!--NOTE: Find a best solution to hide the submit button-->\n      <button style=\"display:none\" type=\"submit\" [disabled]=\"!elemForm.valid\">\n\n      </button>\n      <button type=\"button\"  (click)=\"finishSendList(elemForm)\" *ngIf=\"isListAction\" title=\"Finish send\">\n        <span><img src=\"../../assets/images/send.png\" alt=\"send\"></span>\n      </button>\n    </form>\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.conversation/conversation.component.main.scss":
/***/ (function(module, exports) {

module.exports = "#oty-conversation {\n  position: relative;\n  height: 100%;\n  background-color: #FFFFFF;\n  border-right: 1px solid #B6B6B6; }\n  #oty-conversation.oty-with-input oty-conversation-history {\n    height: calc(100% - 60px); }\n  #oty-conversation.oty-with-input .oty-conversation-form {\n    display: block; }\n  #oty-conversation oty-conversation-history {\n    display: block;\n    height: 100%; }\n  #oty-conversation .oty-conversation-form {\n    display: none;\n    position: absolute;\n    left: 24px;\n    bottom: 24px;\n    width: calc(100% - 48px); }\n  #oty-conversation .oty-conversation-form input {\n      float: left;\n      display: inline-block;\n      width: calc(100% - 12px);\n      height: 33px;\n      border: 1px solid #E0F0F9;\n      color: #212121;\n      border-radius: 4px;\n      padding: 0 12px;\n      font-size: 13px; }\n  #oty-conversation .oty-conversation-form input:focus {\n        outline: none; }\n  #oty-conversation .oty-conversation-form button {\n      float: left;\n      display: inline-block;\n      background-color: #2694D3;\n      border: 0;\n      padding: 4px;\n      width: 24px;\n      height: 24px;\n      border-radius: 50%;\n      margin-top: 5px;\n      margin-left: 10px; }\n  #oty-conversation .oty-conversation-form button:focus {\n        outline: none; }\n  #oty-conversation .oty-conversation-form button:hover {\n        cursor: pointer;\n        background-color: #0081cb; }\n  #oty-conversation .oty-conversation-form button span {\n        color: white; }\n  #oty-conversation .oty-conversation-form button img {\n        width: 13px;\n        height: 13px;\n        vertical-align: 0; }\n  #oty-conversation .oty-conversation-form.is-list input {\n      width: calc(100% - 40px); }\n  #oty-conversation .oty-conversation-form.is-list button {\n      right: 50px; }\n  #oty-conversation .oty-conversation-form.is-list button[type=\"button\"] {\n      right: 8px; }\n"

/***/ }),

/***/ "./src/app/ontology.conversation/conversation.component.ts":
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
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const forms_1 = __webpack_require__("./node_modules/@angular/forms/esm2015/forms.js");
const conversation_service_1 = __webpack_require__("./src/app/services/conversation.service.ts");
const conversation_history_component_1 = __webpack_require__("./src/app/ontology.conversation.history/conversation-history.component.ts");
__webpack_require__("./node_modules/rxjs/_esm2015/add/operator/catch.js");
const notification_service_1 = __webpack_require__("./src/app/services/notification.service.ts");
const cytoscape_service_1 = __webpack_require__("./src/app/services/cytoscape.service.ts");
const PQueue = __webpack_require__("./node_modules/p-queue/index.js");
const layout_service_1 = __webpack_require__("./src/app/ontology.layout/layout.service.ts");
let ConversationComponent = class ConversationComponent {
    constructor(cyService, notificationService, conversationService, router, layoutService, fb) {
        this.cyService = cyService;
        this.notificationService = notificationService;
        this.conversationService = conversationService;
        this.router = router;
        this.layoutService = layoutService;
        this.isListAction = false;
        this.buttonActions = [];
        this.multiselectOptions = [];
        this.formBuilder = fb;
        this.router.params.subscribe(params => {
            this.projectId = params['id'];
        });
    }
    ngOnInit() {
        this.conversationService.clearConversationHistoryList();
        this.conversationService.startConversation(this.projectId).subscribe((res) => {
            if (res) {
                this.conversationService.initConversationHistoryList(res.botResponseData.value.conversationHistory);
                this.formatConversation(res);
            }
        });
        this.conversationService.onRerunCurrentConversation.subscribe((unfinishedTopic) => {
            this.conversationService.resumeTopic(this.projectId, unfinishedTopic.id).subscribe((res) => {
                this.formatConversation(res);
                this.paintGraphPanel(res);
                this.paintRightSidePanel(res);
            });
        });
        this.conversationService.onFormatConversation.subscribe((res) => {
            this.formatConversation(res);
        });
        this.elemForm = this.formBuilder.group({
            'userInput': ["", [forms_1.Validators.required, this.isEmptyString, this.isValidList.bind(this)]],
        });
    }
    continueConversation(userRequest) {
        this.conversationService.continueConversation(this.projectId, userRequest).subscribe((res) => {
            this.formatConversation(res);
            this.paintGraphPanel(res);
            this.paintRightSidePanel(res);
        });
    }
    getConversationHistoryList() {
        return this.conversationService.getConversationHistoryList();
    }
    /**
     * Set bot questions, set button options, and set chat history.
     * @param {ConversationRes} res
     */
    formatConversation(res) {
        if (res) {
            this.conversationService.setDialogHistory(res.dialogHistory);
            this.conversationService.addConversationHistoryList(res);
            //Note: reset below values to default data
            this.resetFormatConversationData();
            let botUserInput = res.botResponseData.value.userInput;
            if (botUserInput.selectableActions) {
                this.buttonActions = botUserInput.selectableActions.actions;
                this.buttonActionSchema = botUserInput.selectableActions.schema;
            }
            if (botUserInput.multiselect) {
                botUserInput.multiselect.options.forEach((text, i) => {
                    this.multiselectOptions.push({
                        text: text,
                        value: i,
                        checked: false
                    });
                });
                this.multiselectSchema = botUserInput.multiselect.schema;
            }
            if (botUserInput.text) {
                this.inputSchema = botUserInput.text.schema;
                this.inputPlaceholder = botUserInput.text.desc;
            }
            if (botUserInput.list) {
                this.inputSchema = botUserInput.list.schema;
                this.inputPlaceholder = botUserInput.list.desc;
                this.isListAction = true;
            }
        }
    }
    resetFormatConversationData() {
        this.buttonActions = [];
        this.buttonActionSchema = undefined;
        this.multiselectOptions = [];
        this.multiselectSchema = undefined;
        this.inputSchema = undefined;
        this.inputPlaceholder = "";
        this.isListAction = false;
    }
    finishSendList(form) {
        this.sendTyping(form, true, true);
    }
    //NOTE: This should be return null if is valid or return a error message if is invalid
    //https://angular.io/guide/form-validation
    isEmptyString(control) {
        let isEmpty = null;
        if ((control.value || '').trim().length === 0) {
            isEmpty = { 'whitespace': true };
        }
        return isEmpty;
    }
    isValidList(control) {
        let userInput = control.value;
        let isValid = true;
        if ((userInput || '').trim().length > 0) {
            userInput = this.filterList(userInput);
            if (userInput === '') {
                isValid = { 'whitespace': true };
            }
        }
        return isValid;
    }
    filterList(value) {
        return value
            .split(',')
            .filter(item => item.trim() !== '')
            .map(item => item.trim())
            .join(',');
    }
    sendTyping(form, isList, isLastItem = false) {
        if (this.inputSchema !== undefined && Object.keys(this.inputSchema).length) {
            const userInput = form.value.userInput;
            let userRequest = {
                key: this.inputSchema.key,
                type: this.inputSchema.type,
                value: userInput
            };
            if (isList) {
                userRequest.value = {
                    item: userInput,
                    lastItem: isLastItem
                };
            }
            this.continueConversation(userRequest);
            form.reset();
        }
    }
    sendButtonOption(value) {
        if (this.buttonActionSchema !== undefined && Object.keys(this.buttonActionSchema).length) {
            let userRequest = {
                key: this.buttonActionSchema.key,
                type: this.buttonActionSchema.type,
                value: value
            };
            this.continueConversation(userRequest);
        }
    }
    sendMultiOptions(value) {
        if (this.multiselectSchema !== undefined && Object.keys(this.multiselectSchema).length) {
            let userRequest = {
                key: this.multiselectSchema.key,
                type: this.multiselectSchema.type,
                value: value
            };
            this.continueConversation(userRequest);
        }
    }
    paintGraphPanel(res) {
        if (res && res.botResponseData.value.commands) {
            let item = res.botResponseData.value.commands.succeed[0];
            this.cyService.sourceType = item.arguments._sourceType;
            console.log("[conversation.component::paintGraphPanel]---(entityName:" + item.arguments.entityName + ")-->item.arguments._sourceType=" + item.arguments._sourceType);
            //Note: expand all the node before save it, otherwise can not get cytoscape json string.
            this.cyService.expandAll();
            this.runCommandsInOrder(res.botResponseData.value.commands.succeed).then(() => {
                this.cyService.cy.center();
                this.conversationService.saveGraphData(this.projectId, this.cyService.cy.json()).subscribe();
            });
        }
    }
    paintRightSidePanel(res) {
        if (res && res.botResponseData.value.graph) {
            this.layoutService.setDetailData(res.botResponseData.value.graph);
        }
    }
    runCommandsInOrder(commands) {
        const queue = new PQueue({ concurrency: 1 });
        const commandsPromise = commands.map((item) => {
            return () => {
                return this[item.command](item.arguments);
            };
        });
        return queue.addAll(commandsPromise);
    }
    getRandomInt(min, max) {
        return Math.floor(Math.random() * (max - min + 1)) + min;
    }
    createEntity(obj) {
        return new Promise((resolve) => {
            let width = this.cyService.cy.width(), height = this.cyService.cy.height();
            this.cyService.cy.add({
                group: "nodes",
                data: {
                    id: obj.entityId,
                    name: obj.entityName
                },
                renderedPosition: {
                    x: this.getRandomInt(width / 4, width * 3 / 4),
                    y: this.getRandomInt(height / 4, height * 3 / 4)
                }
            });
            resolve();
        });
    }
    renameEntity(obj) {
        return new Promise((resolve) => {
            this.cyService.cy.$('#' + obj.entityId).data("name", obj.newName);
            resolve();
        });
    }
    deleteEntity(obj) {
        return new Promise((resolve) => {
            this.cyService.cy.$('#' + obj.entityId).remove();
            resolve();
        });
    }
    createRelation(obj) {
        return new Promise((resolve) => {
            this.cyService.cy.add({
                group: "edges",
                data: {
                    id: obj.relationId,
                    target: obj.targetEntityId,
                    name: obj.relationName,
                    source: obj.sourceEntityId
                },
            });
            resolve();
        });
    }
    renameRelation(obj) {
        return new Promise((resolve) => {
            this.cyService.cy.$('#' + obj.relationId).data("name", obj.newName);
            resolve();
        });
    }
    deleteRelation(obj) {
        return new Promise((resolve) => {
            this.cyService.cy.$('#' + obj.relationId).remove();
            resolve();
        });
    }
    createAttribute(obj) {
        return new Promise((resolve) => {
            let width = this.cyService.cy.width(), height = this.cyService.cy.height();
            //Note: parent container node, child attribute node is not linkable
            this.cyService.cy.add({
                group: "nodes",
                data: {
                    id: obj.attributeId,
                    name: obj.attributeLabel,
                    parent: obj.entityId
                },
                renderedPosition: {
                    x: this.getRandomInt(width / 4, width / 2),
                    y: this.getRandomInt(height / 4, height / 2)
                },
                classes: "attribute"
            });
            resolve();
        });
    }
    deleteAttribute(obj) {
        return new Promise((resolve) => {
            this.cyService.cy.$('#' + obj.attributeId).remove();
            resolve();
        });
    }
    renameAttribute(obj) {
        return new Promise((resolve) => {
            this.cyService.cy.$('#' + obj.attributeId).data("name", obj.newName);
            resolve();
        });
    }
};
__decorate([
    core_1.ViewChild(conversation_history_component_1.ConversationHistoryComponent),
    __metadata("design:type", conversation_history_component_1.ConversationHistoryComponent)
], ConversationComponent.prototype, "historyComponent", void 0);
ConversationComponent = __decorate([
    core_1.Component({
        selector: 'oty-conversation',
        template: __webpack_require__("./src/app/ontology.conversation/conversation.component.html"),
        styles: [__webpack_require__("./src/app/ontology.conversation/conversation.component.main.scss")]
    }),
    __param(5, core_1.Inject(forms_1.FormBuilder)),
    __metadata("design:paramtypes", [cytoscape_service_1.CytoscapeService,
        notification_service_1.NotificationService,
        conversation_service_1.ConversationService,
        router_1.ActivatedRoute,
        layout_service_1.LayoutService,
        forms_1.FormBuilder])
], ConversationComponent);
exports.ConversationComponent = ConversationComponent;


/***/ }),

/***/ "./src/app/ontology.conversation/conversation.module.ts":
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
const conversation_component_1 = __webpack_require__("./src/app/ontology.conversation/conversation.component.ts");
const api_1 = __webpack_require__("./node_modules/primeng/api.js");
const shared_module_1 = __webpack_require__("./src/app/common/shared.module.ts");
const conversation_history_module_1 = __webpack_require__("./src/app/ontology.conversation.history/conversation-history.module.ts");
let ConversationModule = class ConversationModule {
};
ConversationModule = __decorate([
    core_1.NgModule({
        imports: [common_1.CommonModule, shared_module_1.SharedModule, conversation_history_module_1.ConversationHistoryModule],
        exports: [conversation_component_1.ConversationComponent],
        declarations: [conversation_component_1.ConversationComponent],
        providers: [api_1.TreeDragDropService]
    })
], ConversationModule);
exports.ConversationModule = ConversationModule;


/***/ }),

/***/ "./src/app/ontology.layout.detail/layout.detail.component.html":
/***/ (function(module, exports) {

module.exports = "<div id=\"oty-layout-detail\">\n  <div class=\"oty-entity\" *ngIf=\"(isNode(selectedElem) || isAttribute(selectedElem)) && graphEntity\">\n    <div class=\"oty-header\">\n      <h4><i class=\"oty-title-circle\"></i><span>{{graphEntity.name}}</span></h4>\n\n      <span class=\"oty-icon-trashcan\" (click)=\"confirmDeleteEntity(graphEntity.id, graphEntity.name)\"></span>\n      <i class=\"oty-status\" [ngClass]=\"{'oty-icon-complete': graphEntity.status==='Complete', 'oty-icon-incomplete': graphEntity.status==='Incomplete'}\"\n         aria-hidden=\"true\"></i>\n    </div>\n\n    <table class=\"oty-table-definition oty-table\" *ngIf=\"graphEntity.definition?.length\">\n      <caption>Definition:</caption>\n      <tr *ngFor=\"let definition of graphEntity.definition; index as i;\">\n        <td width=\"85%\">{{definition}}</td>\n        <td width=\"15%\">\n          <div class=\"oty-table-delete-button\">\n            <i class=\"oty-icon-remove\" aria-hidden=\"true\"></i>\n            <a (click)=\"confirmDeleteDefinition(graphEntity.name, graphEntity.definition, i)\">Yes</a>\n          </div>\n        </td>\n      </tr>\n    </table>\n\n    <div class=\"oty-entity-keys\">\n      <h5>Key: </h5>\n      <span>{{graphEntity.key && graphEntity.key.length ? graphEntity.key.join(\", \") : 'undefined'}}</span>\n    </div>\n\n        <table class=\"oty-table-attribute oty-table\" *ngIf=\"graphEntity.attributes?.length\">\n          <caption>Attributes:</caption>\n          <tr>\n            <td width=\"20%\">Name</td>\n            <td width=\"40%\">value</td>\n            <td width=\"20%\">type</td>\n            <td width=\"20%\">Remove</td>\n          </tr>\n          <tr *ngFor=\"let attribute of graphEntity.attributes; index as i;\">\n            <td>{{attribute.name}}</td>\n            <td>{{attribute.value}}</td>\n            <td>{{attribute.type}}</td>\n            <td>\n              <div class=\"oty-table-delete-button\">\n                <i class=\"oty-icon-remove\" aria-hidden=\"true\"></i>\n                <a (click)=\"confirmDeleteAttribute(graphEntity.name, attribute.id, attribute.name)\">Yes</a>\n              </div>\n            </td>\n          </tr>\n        </table>\n\n    <table class=\"oty-table-relation oty-table\" *ngIf=\"graphEntity.relations?.length\">\n      <caption>Relationships to other entities:</caption>\n      <tr>\n        <td width=\"30%\">Target name</td>\n        <td width=\"40%\">Source name</td>\n        <td width=\"15%\">Name</td>\n        <td width=\"15%\">Remove</td>\n      </tr>\n      <tr *ngFor=\"let relation of graphEntity.relations; index as i;\">\n        <td><span (click)=\"selectElement(relation.targetEntityId)\">{{relation.targetEntityName}}</span></td>\n        <td><span (click)=\"selectElement(relation.sourceEntityId)\">{{relation.sourceEntityName}}</span></td>\n        <td>{{relation.name}}</td>\n        <td>\n          <div class=\"oty-table-delete-button\">\n            <i class=\"oty-icon-remove\" aria-hidden=\"true\"></i>\n            <a (click)=\"confirmDeleteRelation(relation.id, relation.name, relation.sourceEntityName, relation.targetEntityName, true)\">Yes</a>\n          </div>\n        </td>\n      </tr>\n    </table>\n\n    <div class=\"oty-no-connection\" *ngIf=\"graphEntity.relations?.length === 0 && graphEntity.attributes?.length === 0\">\n      <h5>Connections</h5>\n      <span>You don't have any connections yet</span>\n    </div>\n  </div>\n\n  <div class=\"oty-relation\" *ngIf=\"isEdge(selectedElem) && graphRelation\">\n    <div class=\"oty-header\">\n      <h4><i class=\"oty-title-circle\"></i><span>{{graphRelation.name || \"Undefined\"}}</span></h4>\n\n      <span class=\"oty-icon-trashcan\" (click)=\"confirmDeleteRelation(graphRelation.id, graphRelation.name, graphRelation.sourceEntityName, graphRelation.targetEntityName, false)\"></span>\n      <i class=\"oty-status\" [ngClass]=\"{'oty-icon-complete': graphRelation.status==='Complete', 'oty-icon-incomplete': graphRelation.status==='Incomplete'}\"\n         aria-hidden=\"true\"></i>\n    </div>\n    <table class=\"oty-table\">\n      <caption>Entities:</caption>\n      <tr>\n        <td width=\"35%\">Target name</td>\n        <td width=\"35%\">Source name</td>\n        <td width=\"30%\">Name</td>\n      </tr>\n      <tr>\n        <td><span (click)=\"selectElement(graphRelation.targetEntityId)\">{{graphRelation.targetEntityName}}</span></td>\n        <td><span (click)=\"selectElement(graphRelation.sourceEntityId)\">{{graphRelation.sourceEntityName}}</span></td>\n        <td>{{graphRelation.name}}</td>\n      </tr>\n    </table>\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.layout.detail/layout.detail.component.main.scss":
/***/ (function(module, exports) {

module.exports = "#oty-layout-detail {\n  float: right;\n  font-size: 13px;\n  position: relative;\n  height: 100%;\n  width: 100%;\n  border-left: 1px solid #B6B6B6;\n  overflow: auto;\n  z-index: 2;\n  padding: 20px;\n  background-color: #FFFFFF; }\n  #oty-layout-detail .oty-header {\n    font-size: 16px;\n    font-weight: 600;\n    position: relative; }\n  #oty-layout-detail .oty-header h4 {\n      margin: 0;\n      height: 20px; }\n  #oty-layout-detail .oty-header h4 .oty-title-circle {\n        border-radius: 50%;\n        height: 16px;\n        width: 16px;\n        border: 4px solid #E0F0F9;\n        background-color: #2694D3;\n        margin-right: 12px;\n        display: inline-block;\n        margin-top: 2px;\n        float: left; }\n  #oty-layout-detail .oty-header .oty-icon-trashcan {\n      position: absolute;\n      top: 3px;\n      right: 25px;\n      cursor: pointer; }\n  #oty-layout-detail .oty-header .oty-icon-trashcan:hover {\n        color: #D0021B; }\n  #oty-layout-detail .oty-header .oty-status {\n      position: absolute;\n      top: 3px;\n      right: 0;\n      display: block; }\n  #oty-layout-detail .oty-entity-keys {\n    margin: 20px 0; }\n  #oty-layout-detail .oty-entity-keys h5 {\n      margin: 0 0 5px 0; }\n  #oty-layout-detail .oty-icon-complete, #oty-layout-detail .oty-icon-incomplete {\n    color: #FFFFFF;\n    font-size: 17px; }\n  #oty-layout-detail .oty-icon-complete {\n    color: #04BE5B; }\n  #oty-layout-detail .oty-icon-incomplete {\n    color: #D0021B; }\n  #oty-layout-detail h5 {\n    line-height: 19px;\n    font-weight: 600;\n    text-align: left; }\n  #oty-layout-detail .oty-table {\n    margin-top: 20px;\n    table-layout: fixed;\n    width: 100%; }\n  #oty-layout-detail .oty-table caption {\n      text-align: left;\n      border-bottom: 1px solid #DEDEDE;\n      font-weight: 600; }\n  #oty-layout-detail .oty-table tr {\n      border-bottom: 1px solid #DEDEDE; }\n  #oty-layout-detail .oty-table tr td {\n        padding: 5px; }\n  #oty-layout-detail .oty-table tr td span {\n          color: #0081cb;\n          cursor: pointer; }\n  #oty-layout-detail .oty-table tr td .oty-table-delete-button {\n          display: inline-block;\n          width: 100%;\n          height: 100%; }\n  #oty-layout-detail .oty-table tr td .oty-table-delete-button i {\n            color: #FFFFFF;\n            background-color: #DEDEDE;\n            font-size: 20px;\n            float: left;\n            width: 19px; }\n  #oty-layout-detail .oty-table tr td .oty-table-delete-button a {\n            float: left;\n            cursor: pointer;\n            color: #0081cb;\n            display: none;\n            margin-left: 5px;\n            margin-top: 2px; }\n  #oty-layout-detail .oty-table tr td .oty-table-delete-button:hover i {\n            background-color: #FF0000; }\n  #oty-layout-detail .oty-table tr td .oty-table-delete-button:hover a {\n            display: block;\n            text-decoration: none; }\n  #oty-layout-detail .oty-table.oty-table-attribute {\n    border-bottom: 1px solid #DEDEDE; }\n  #oty-layout-detail .oty-table.oty-table-attribute caption, #oty-layout-detail .oty-table.oty-table-attribute tr:first-child {\n      border-bottom: 1px solid #DEDEDE; }\n  #oty-layout-detail .oty-table.oty-table-attribute tr {\n      border-bottom: 0; }\n  #oty-layout-detail .oty-table.oty-table-attribute tr:nth-child(n+1) td:nth-child(3n) {\n        text-align: center; }\n  #oty-layout-detail .oty-no-connection, #oty-layout-detail .oty-specification {\n    margin-top: 20px; }\n"

/***/ }),

/***/ "./src/app/ontology.layout.detail/layout.detail.component.ts":
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
const cytoscape_service_1 = __webpack_require__("./src/app/services/cytoscape.service.ts");
const project_service_1 = __webpack_require__("./src/app/services/project.service.ts");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const conversation_service_1 = __webpack_require__("./src/app/services/conversation.service.ts");
const layout_service_1 = __webpack_require__("./src/app/ontology.layout/layout.service.ts");
let LayoutDetailComponent = class LayoutDetailComponent {
    constructor(cyService, projectService, conversationService, layoutService, router) {
        this.cyService = cyService;
        this.projectService = projectService;
        this.conversationService = conversationService;
        this.layoutService = layoutService;
        this.router = router;
        this.subscriptions = [];
        this.router.params.subscribe(params => {
            this.projectId = params['id'];
        });
        this.selectedElem = null;
    }
    ngOnInit() {
        this.subscriptions.push(this.cyService.event.subscribe(event => {
            if (typeof this[event.name] === "function") {
                this[event.name](event.event);
            }
        }));
        this.subscriptions.push(this.layoutService.onSetDetailData.subscribe((graph) => {
            this.setDetailData(graph);
        }));
    }
    ngOnDestroy() {
        this.subscriptions.forEach((subscription) => subscription.unsubscribe());
    }
    setDetailData(graph) {
        this.graphEntity = undefined;
        this.graphRelation = undefined;
        if (graph && graph.entity) {
            this.graphEntity = this.preFormatGraphEntity(graph.entity);
        }
        if (graph && graph.relation) {
            this.graphRelation = graph.relation;
        }
    }
    preFormatGraphEntity(entity) {
        entity.key = entity.key.filter(item => item);
        return entity;
    }
    onNothingClick(event) {
        this.layoutService.toggleDetailPanel(false, true);
        this.setSelectElem(null);
        this.graphEntity = undefined;
        this.graphRelation = undefined;
        this.cyService.collapseAll();
    }
    onNodeClick(event) {
        this.setSelectElem(event.target);
        this.layoutService.toggleDetailPanel(true, true);
        let entityName = event.target.data("name");
        this.conversationService.getEntityDetail(this.projectId, entityName).subscribe((res) => {
            this.conversationService.fireEventFormatConversation(res);
            this.setDetailData(res.botResponseData.value.graph);
        });
    }
    onAttributeClick(event) {
        this.setSelectElem(event.target);
        this.layoutService.toggleDetailPanel(true, true);
        let entityName = event.target.parent().data("name");
        let attributeName = event.target.data("name");
        let attributeId = event.target.data("id");
        this.conversationService.getAttributeDetail(this.projectId, entityName, attributeId).subscribe((res) => {
            this.conversationService.fireEventFormatConversation(res);
            this.setDetailData(res.botResponseData.value.graph);
        });
    }
    onEdgeClick(event) {
        this.setSelectElem(event.target);
        this.layoutService.toggleDetailPanel(true, true);
        let relationName = event.target.data("name");
        let sourceEntityName = event.target.source().data("name");
        let targetEntityName = event.target.target().data("name");
        this.conversationService.getRelationDetail(this.projectId, relationName, sourceEntityName, targetEntityName).subscribe((res) => {
            this.conversationService.fireEventFormatConversation(res);
            this.setDetailData(res.botResponseData.value.graph);
        });
    }
    onEdgeCreate(event) {
        this.setSelectElem(event.data.addedEles);
        this.layoutService.toggleDetailPanel(true, true);
    }
    setSelectElem(data) {
        this.selectedElem = data;
    }
    isNode(ele) {
        return !this.isAttribute(ele);
    }
    isEdge(ele) {
        return ele && ele.isEdge();
    }
    isAttribute(ele) {
        return ele && ele.isNode() && ele.isChild();
    }
    confirmDeleteAttribute(entityName, attributeId, attributeName) {
        this.conversationService.deleteAttribute(this.projectId, entityName, attributeName).subscribe((res) => {
            // Note: Notice conversation list to update
            this.conversationService.fireEventFormatConversation(res);
            let commandObj = this.getDeleteCommand(res);
            if (commandObj && commandObj.command === "deleteAttribute") {
                this.layoutService.toggleDetailPanel(false, true);
                //Note: expand the node before delete it, otherwise can not get cytoscape json string.
                this.cyService.expandAll();
                this.cyService.deleteElement(commandObj.arguments.attributeId);
                this.graphEntity.attributes = this.graphEntity.attributes.filter((attribute) => {
                    return attribute.id !== commandObj.arguments.attributeId;
                });
                this.saveGraphData();
            }
        });
    }
    confirmDeleteEntity(entityId, entityName) {
        this.conversationService.deleteEntity(this.projectId, entityName).subscribe((res) => {
            // Note: Notice conversation list to update
            this.conversationService.fireEventFormatConversation(res);
            let commandObj = this.getDeleteCommand(res);
            if (commandObj && commandObj.command === "deleteEntity") {
                this.layoutService.toggleDetailPanel(false, true);
                //Note: expand the node before delete it, otherwise can not get cytoscape json string.
                this.cyService.expandAll();
                this.cyService.deleteElement(commandObj.arguments.entityId);
                this.saveGraphData();
            }
        });
    }
    confirmDeleteDefinition(entityName, definitions, deleteIndex) {
        this.conversationService.deleteDefinition(this.projectId, entityName, deleteIndex).subscribe((res) => {
            // Note: Notice conversation list to update
            this.conversationService.fireEventFormatConversation(res);
        });
    }
    confirmDeleteRelation(relationId, relationName, sourceEntityName, targetEntityName, isShowRightPanel) {
        this.conversationService.deleteRelation(this.projectId, relationName, sourceEntityName, targetEntityName).subscribe((res) => {
            // Note: Notice conversation list to update
            this.conversationService.fireEventFormatConversation(res);
            let commandObj = this.getDeleteCommand(res);
            if (commandObj && commandObj.command === "deleteRelation") {
                //Note: expand the node before delete it, otherwise can not get cytoscape json string.
                this.cyService.expandAll();
                this.layoutService.toggleDetailPanel(isShowRightPanel, true);
                this.cyService.deleteElement(commandObj.arguments.relationId);
                this.graphEntity.relations = this.graphEntity.relations.filter((relation) => {
                    return relation.id !== commandObj.arguments.relationId;
                });
                this.saveGraphData();
            }
        });
    }
    getDeleteCommand(res) {
        let command;
        if (res.botResponseData.value.commands && res.botResponseData.value.commands.succeed.length) {
            command = res.botResponseData.value.commands.succeed[0];
        }
        return command;
    }
    saveGraphData() {
        this.conversationService.saveGraphData(this.projectId, this.cyService.cy.json()).subscribe();
    }
    selectElement(elemId) {
        this.cyService.selectElement(elemId);
    }
};
LayoutDetailComponent = __decorate([
    core_1.Component({
        selector: 'oty-layout-detail',
        template: __webpack_require__("./src/app/ontology.layout.detail/layout.detail.component.html"),
        styles: [__webpack_require__("./src/app/ontology.layout.detail/layout.detail.component.main.scss")]
    }),
    __metadata("design:paramtypes", [cytoscape_service_1.CytoscapeService,
        project_service_1.ProjectService,
        conversation_service_1.ConversationService,
        layout_service_1.LayoutService,
        router_1.ActivatedRoute])
], LayoutDetailComponent);
exports.LayoutDetailComponent = LayoutDetailComponent;


/***/ }),

/***/ "./src/app/ontology.layout.detail/layout.detail.module.ts":
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
const shared_module_1 = __webpack_require__("./src/app/common/shared.module.ts");
const dragdrop_1 = __webpack_require__("./node_modules/primeng/dragdrop.js");
const inputtext_1 = __webpack_require__("./node_modules/primeng/inputtext.js");
const button_1 = __webpack_require__("./node_modules/primeng/button.js");
const radiobutton_1 = __webpack_require__("./node_modules/primeng/radiobutton.js");
const layout_detail_component_1 = __webpack_require__("./src/app/ontology.layout.detail/layout.detail.component.ts");
const common_1 = __webpack_require__("./node_modules/@angular/common/esm2015/common.js");
let LayoutDetailModule = class LayoutDetailModule {
};
LayoutDetailModule = __decorate([
    core_1.NgModule({
        imports: [common_1.CommonModule, shared_module_1.SharedModule, dragdrop_1.DragDropModule, button_1.ButtonModule, inputtext_1.InputTextModule, radiobutton_1.RadioButtonModule],
        exports: [layout_detail_component_1.LayoutDetailComponent],
        declarations: [layout_detail_component_1.LayoutDetailComponent],
        providers: []
    })
], LayoutDetailModule);
exports.LayoutDetailModule = LayoutDetailModule;


/***/ }),

/***/ "./src/app/ontology.layout.options/layout.options.component.html":
/***/ (function(module, exports) {

module.exports = "<div id=\"oty-option-panel\">\n  <ul>\n    <li><a title=\"Zoom in\" (click)=\"zoomIn()\"><i class=\"oty-icon-zoom-in\" aria-hidden=\"true\"></i></a></li>\n    <li><a title=\"Zoom out\" (click)=\"zoomOut()\"><i class=\"oty-icon-zoom-out\" aria-hidden=\"true\"></i></a></li>\n    <li class=\"dropdown\">\n      <a title=\"Change layout\" href=\"javascript:void(0)\" class=\"dropbtn\"><i class=\"oty-icon-force-link\" aria-hidden=\"true\"></i></a>\n      <div class=\"dropdown-content\" >\n        <a *ngFor=\"let option of options\" (click)=\"updateLayout(option)\">{{option}}</a>\n      </div>\n    </li>\n    <li><a title=\"Reset layout\" (click)=\"centerLayout()\"><i class=\"oty-icon-grid\" aria-hidden=\"true\"></i></a></li>\n    <li><a title=\"Expand nodes\" (click)=\"expandNodes()\"><i class=\"oty-icon-expand\"></i></a></li>\n    <li><a title=\"Collapse nodes\" (click)=\"collapseNodes()\"><i class=\"oty-icon-collapse\"></i></a></li>\n  </ul>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.layout.options/layout.options.component.scss":
/***/ (function(module, exports) {

module.exports = "ul {\n  list-style-type: none;\n  margin: 0;\n  overflow: hidden;\n  padding: 0; }\n  ul .dropdown-content {\n    background-color: #FAFAFA;\n    -webkit-box-shadow: 0px 8px 16px 0px #727272;\n            box-shadow: 0px 8px 16px 0px #727272;\n    display: none;\n    min-width: 160px;\n    position: absolute;\n    right: 100%;\n    top: 30%;\n    z-index: 1; }\n  ul .dropdown-content a {\n      color: #000000;\n      cursor: pointer;\n      display: block;\n      padding: 12px 16px;\n      text-align: left; }\n  ul .dropdown-content a:hover {\n      background-color: #FFFFFF;\n      color: #212121;\n      text-decoration: none; }\n  ul li a,\n  ul .dropbtn {\n    color: #727272;\n    display: inline-block;\n    padding: 10px; }\n  ul li a:hover,\n  ul .dropdown:hover .dropbtn {\n    background-color: #727272;\n    color: #FAFAFA;\n    text-decoration: none; }\n  ul .dropdown {\n    display: inline-block; }\n  ul .dropdown:hover .dropdown-content {\n    display: block;\n    text-decoration: none; }\n  ul [class^=oty-icon-] {\n    color: inherit;\n    font-size: 20px; }\n"

/***/ }),

/***/ "./src/app/ontology.layout.options/layout.options.component.ts":
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
let LayoutOptionsComponent = class LayoutOptionsComponent {
    constructor() {
        this.options = ['grid', 'circle', 'cose-bilkent', 'concentric', 'breadthfirst'];
        this.changeLayout = new core_1.EventEmitter();
        this.zoom = new core_1.EventEmitter();
        this.center = new core_1.EventEmitter();
        this.onExpandNodes = new core_1.EventEmitter();
        this.onCollapseNodes = new core_1.EventEmitter();
    }
    ngOnInit() {
    }
    updateLayout(option) {
        this.changeLayout.emit(option);
    }
    zoomIn() {
        this.zoom.emit(0.1);
    }
    zoomOut() {
        this.zoom.emit(-0.1);
    }
    centerLayout() {
        this.center.emit();
    }
    expandNodes() {
        this.onExpandNodes.emit();
    }
    collapseNodes() {
        this.onCollapseNodes.emit();
    }
};
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], LayoutOptionsComponent.prototype, "changeLayout", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], LayoutOptionsComponent.prototype, "zoom", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], LayoutOptionsComponent.prototype, "center", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], LayoutOptionsComponent.prototype, "onExpandNodes", void 0);
__decorate([
    core_1.Output(),
    __metadata("design:type", core_1.EventEmitter)
], LayoutOptionsComponent.prototype, "onCollapseNodes", void 0);
LayoutOptionsComponent = __decorate([
    core_1.Component({
        selector: 'oty-layout-options',
        template: __webpack_require__("./src/app/ontology.layout.options/layout.options.component.html"),
        styles: [__webpack_require__("./src/app/ontology.layout.options/layout.options.component.scss")]
    }),
    __metadata("design:paramtypes", [])
], LayoutOptionsComponent);
exports.LayoutOptionsComponent = LayoutOptionsComponent;


/***/ }),

/***/ "./src/app/ontology.layout.options/layout.options.module.ts":
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
const layout_options_component_1 = __webpack_require__("./src/app/ontology.layout.options/layout.options.component.ts");
let LayoutOptionsModule = class LayoutOptionsModule {
};
LayoutOptionsModule = __decorate([
    core_1.NgModule({
        imports: [
            common_1.CommonModule
        ],
        exports: [layout_options_component_1.LayoutOptionsComponent],
        declarations: [layout_options_component_1.LayoutOptionsComponent]
    })
], LayoutOptionsModule);
exports.LayoutOptionsModule = LayoutOptionsModule;


/***/ }),

/***/ "./src/app/ontology.layout/layout.component.html":
/***/ (function(module, exports) {

module.exports = "<div id=\"oty-cytoscape-layout-container\" [class.showDetailPanel]=\"isShowSidePanel\">\n  <div #cy id=\"oty-cytoscape-layout\"></div>\n  <oty-layout-options\n    (center)=\"center()\"\n    (changeLayout)=\"updateLayout($event)\"\n    (onCollapseNodes)=\"collapseNodes()\"\n    (onExpandNodes)=\"expandNodes()\"\n    (zoom)=\"zoom($event)\"></oty-layout-options>\n  <oty-layout-detail [hidden]=\"isShowSidePanel\"></oty-layout-detail>\n  <div class=\"oty-selector\" [hidden]=\"!isShowSideIcon\" (click)=\"toggleDetailPanel(!isShowSidePanel)\">\n    <img src=\"../../assets/images/open-tab.png\">\n  </div>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.layout/layout.component.main.scss":
/***/ (function(module, exports) {

module.exports = "#oty-cytoscape-layout-container {\n  position: relative;\n  width: 100%;\n  height: 100%; }\n  #oty-cytoscape-layout-container.showDetailPanel #oty-cytoscape-layout {\n    width: calc(100% - 600px);\n    float: left; }\n  #oty-cytoscape-layout-container.showDetailPanel oty-layout-detail {\n    display: block; }\n  #oty-cytoscape-layout-container.showDetailPanel oty-layout-options {\n    right: 600px; }\n  #oty-cytoscape-layout-container.showDetailPanel .oty-selector {\n    right: 600px; }\n  #oty-cytoscape-layout-container #oty-cytoscape-layout {\n    width: 100%;\n    height: 100%;\n    z-index: 1; }\n  #oty-cytoscape-layout-container oty-layout-options {\n    position: absolute;\n    z-index: 2;\n    width: 34px;\n    top: 0;\n    right: 0; }\n  #oty-cytoscape-layout-container oty-layout-detail {\n    width: 600px;\n    display: none;\n    height: 100%;\n    float: right; }\n  #oty-cytoscape-layout-container .oty-selector {\n    position: absolute;\n    right: 0;\n    top: calc(50% - 34px/2);\n    z-index: 2;\n    cursor: pointer; }\n  #oty-cytoscape-layout-container .oty-selector img {\n      width: 15px; }\n"

/***/ }),

/***/ "./src/app/ontology.layout/layout.component.ts":
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
const cytoscape_service_1 = __webpack_require__("./src/app/services/cytoscape.service.ts");
const router_1 = __webpack_require__("./node_modules/@angular/router/esm2015/router.js");
const project_service_1 = __webpack_require__("./src/app/services/project.service.ts");
const http_1 = __webpack_require__("./node_modules/@angular/common/esm2015/http.js");
const conversation_service_1 = __webpack_require__("./src/app/services/conversation.service.ts");
const layout_service_1 = __webpack_require__("./src/app/ontology.layout/layout.service.ts");
let LayoutComponent = class LayoutComponent {
    constructor(cyService, projectService, conversationService, layoutService, router, http) {
        this.cyService = cyService;
        this.projectService = projectService;
        this.conversationService = conversationService;
        this.layoutService = layoutService;
        this.router = router;
        this.http = http;
        this.subscriptions = [];
        this.router.params.subscribe(params => {
            this.projectId = params['id'];
        });
        this.isShowSidePanel = false;
        this.isShowSideIcon = false;
    }
    ngOnInit() {
        this.initCytoscape();
        this.subscriptions.push(this.conversationService.getGraphData(this.projectId).subscribe((graphData) => {
            if (graphData !== null && graphData.elements !== null) {
                this.cyService.cy.json(graphData);
                this.cyService.expandAll();
            }
        }));
        this.subscriptions.push(this.layoutService.toggleDetailPanelEvent.subscribe((event) => {
            if (event.isRightSide) {
                this.toggleDetailPanel(event.isShow);
            }
        }));
        this.subscriptions.push(this.cyService.event.subscribe(event => {
            if (event.name === "onNothingClick") {
                this.isShowSideIcon = false;
            }
        }));
    }
    ngOnDestroy() {
        this.subscriptions.forEach((subscription) => subscription.unsubscribe());
    }
    initCytoscape() {
        this.cy = this.cyService.initCytoscape(this.cyElement.nativeElement);
    }
    resizeCytoscape() {
        this.cyService.resize();
    }
    updateLayout(option) {
        this.cyService.updateLayout(option);
    }
    zoom(zoomNumber) {
        this.cyService.zoom(zoomNumber);
    }
    collapseNodes() {
        this.cyService.collapseAll();
    }
    expandNodes() {
        this.cyService.expandAll();
    }
    center() {
        this.cyService.center();
    }
    toggleDetailPanel(isShowSidePanel) {
        this.isShowSidePanel = isShowSidePanel;
        if (isShowSidePanel) {
            this.isShowSideIcon = true;
        }
        setTimeout(() => {
            this.resizeCytoscape();
            this.center();
        }, 10);
    }
};
__decorate([
    core_1.ViewChild('cy'),
    __metadata("design:type", core_1.ElementRef)
], LayoutComponent.prototype, "cyElement", void 0);
LayoutComponent = __decorate([
    core_1.Component({
        selector: 'oty-layout',
        template: __webpack_require__("./src/app/ontology.layout/layout.component.html"),
        styles: [__webpack_require__("./src/app/ontology.layout/layout.component.main.scss")]
    }),
    __metadata("design:paramtypes", [cytoscape_service_1.CytoscapeService,
        project_service_1.ProjectService,
        conversation_service_1.ConversationService,
        layout_service_1.LayoutService,
        router_1.ActivatedRoute, http_1.HttpClient])
], LayoutComponent);
exports.LayoutComponent = LayoutComponent;


/***/ }),

/***/ "./src/app/ontology.layout/layout.module.ts":
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
const dragdrop_1 = __webpack_require__("./node_modules/primeng/dragdrop.js");
const layout_component_1 = __webpack_require__("./src/app/ontology.layout/layout.component.ts");
const layout_options_module_1 = __webpack_require__("./src/app/ontology.layout.options/layout.options.module.ts");
const layout_detail_module_1 = __webpack_require__("./src/app/ontology.layout.detail/layout.detail.module.ts");
let LayoutModule = class LayoutModule {
};
LayoutModule = __decorate([
    core_1.NgModule({
        imports: [common_1.CommonModule, forms_1.FormsModule, dragdrop_1.DragDropModule, layout_options_module_1.LayoutOptionsModule, layout_detail_module_1.LayoutDetailModule],
        exports: [layout_component_1.LayoutComponent],
        declarations: [layout_component_1.LayoutComponent],
        providers: []
    })
], LayoutModule);
exports.LayoutModule = LayoutModule;


/***/ }),

/***/ "./src/app/ontology.project/project-routing.module.ts":
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
const project_component_1 = __webpack_require__("./src/app/ontology.project/project.component.ts");
const auth_guard_service_1 = __webpack_require__("./src/app/services/auth-guard.service.ts");
let ProjectRoutingModule = class ProjectRoutingModule {
};
ProjectRoutingModule = __decorate([
    core_1.NgModule({
        imports: [router_1.RouterModule.forChild([
                { path: '', component: project_component_1.ProjectComponent, canActivate: [auth_guard_service_1.AuthGuardService] }
            ])],
        exports: [router_1.RouterModule]
    })
], ProjectRoutingModule);
exports.ProjectRoutingModule = ProjectRoutingModule;


/***/ }),

/***/ "./src/app/ontology.project/project.component.html":
/***/ (function(module, exports) {

module.exports = "<div id=\"oty-header-options\">\n  <a routerLink=\"/home\">\n    <i class=\"oty-icon-back\" aria-hidden=\"true\"></i>\n    <img src=\"../../assets/images/project/project1.svg\">\n    <span>{{projectData?.projectName}}</span>\n  </a>\n  <div>\n    <span *ngIf=\"projectData?.timestamp\">Last modified on {{projectData.timestamp | date:'long'}}</span>\n    <span *ngIf=\"projectData?.username\"> by {{projectData.username}}</span>\n    <a *ngIf=\"projectData?.projectName\" class=\"tc-buttons\" [href]=\"fileUrl\" download=\"{{projectData.projectName + '.json'}}\">\n      <i class=\"tc-icon oty-icon-export\"></i>\n      <span>Export</span>\n    </a>\n  </div>\n</div>\n<div id=\"oty-project\">\n\n  <split direction=\"horizontal\" gutterSize=\"8\">\n    <split-area size=\"30\">\n      <oty-unfinished-topic></oty-unfinished-topic>\n      <oty-conversation></oty-conversation>\n    </split-area>\n    <split-area size=\"70\">\n      <oty-layout></oty-layout>\n    </split-area>\n  </split>\n</div>\n"

/***/ }),

/***/ "./src/app/ontology.project/project.component.main.scss":
/***/ (function(module, exports) {

module.exports = "#oty-header-options {\n  background-color: #FFFFFF;\n  border-bottom: 1px solid #DEDEDE;\n  -webkit-box-sizing: border-box;\n          box-sizing: border-box;\n  -webkit-box-shadow: 0 2px 4px 0 rgba(89, 117, 168, 0.09);\n          box-shadow: 0 2px 4px 0 rgba(89, 117, 168, 0.09);\n  width: 100%;\n  height: 74px; }\n  #oty-header-options > a {\n    display: inline-block;\n    height: 44px;\n    line-height: 44px;\n    margin: 15px 0 0 25px; }\n  #oty-header-options > a:hover {\n      text-decoration: none; }\n  #oty-header-options > a .oty-icon-back {\n      font-size: 24px;\n      color: #0081cb;\n      background-color: #FFFFFF;\n      margin-top: 8px;\n      display: inline-block;\n      float: left; }\n  #oty-header-options > a img {\n      height: 44px;\n      width: 48px;\n      float: left;\n      margin-left: 40px; }\n  #oty-header-options > a span {\n      font-size: 24px;\n      font-weight: 600;\n      letter-spacing: 0.4px;\n      line-height: 44px;\n      margin-left: 20px;\n      float: left; }\n  #oty-header-options > div {\n    display: inline-block;\n    width: 50%;\n    float: right;\n    text-align: right; }\n  #oty-header-options > div a {\n      background-color: #0081cb;\n      border-radius: 5px;\n      color: #FFFFFF;\n      font-size: 16px;\n      font-weight: 600;\n      height: 42px;\n      text-align: center;\n      margin-top: 16px;\n      margin-right: 25.16px;\n      margin-left: 17px;\n      width: 121.68px;\n      padding-left: 20px; }\n  #oty-header-options > div a i {\n        font-size: 24px;\n        float: left; }\n  #oty-header-options > div a:focus {\n        outline: none; }\n  #oty-header-options > div a:hover {\n        text-decoration: none;\n        background-color: #2694D3; }\n  #oty-header-options > div a:active {\n        text-decoration: none; }\n  #oty-header-options > div > span {\n      display: inline-block;\n      position: relative;\n      font-size: 14px;\n      letter-spacing: 0.3px;\n      line-height: 18px; }\n  #oty-project {\n  display: -webkit-box;\n  display: -ms-flexbox;\n  display: flex;\n  height: calc(100% - 74px);\n  position: relative; }\n  #oty-project split-area {\n    position: relative; }\n  #oty-project split-area:first-child {\n      min-width: 300px; }\n  #oty-project oty-conversation {\n    min-width: 300px;\n    width: 300px; }\n  #oty-project oty-layout {\n    width: calc(100% - 300px); }\n"

/***/ }),

/***/ "./src/app/ontology.project/project.component.ts":
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
const platform_browser_1 = __webpack_require__("./node_modules/@angular/platform-browser/esm2015/platform-browser.js");
const project_service_1 = __webpack_require__("./src/app/services/project.service.ts");
let ProjectComponent = class ProjectComponent {
    constructor(projectService, router, sanitizer) {
        this.projectService = projectService;
        this.router = router;
        this.sanitizer = sanitizer;
        this.router.params.subscribe(params => {
            this.id = params['id'];
        });
    }
    ngOnInit() {
        this.projectService.getProject(this.id)
            .subscribe((res) => {
            this.projectData = res;
        });
        this.getExportFileUrl();
    }
    getExportFileUrl() {
        this.projectService.exportProject(this.id).subscribe((res) => {
            const blob = new Blob([JSON.stringify(res)], { type: 'application/octet-stream' });
            this.fileUrl = this.sanitizer.bypassSecurityTrustResourceUrl(window.URL.createObjectURL(blob));
        });
    }
};
ProjectComponent = __decorate([
    core_1.Component({
        selector: 'oty-project',
        template: __webpack_require__("./src/app/ontology.project/project.component.html"),
        styles: [__webpack_require__("./src/app/ontology.project/project.component.main.scss")]
    }),
    __metadata("design:paramtypes", [project_service_1.ProjectService,
        router_1.ActivatedRoute,
        platform_browser_1.DomSanitizer])
], ProjectComponent);
exports.ProjectComponent = ProjectComponent;


/***/ }),

/***/ "./src/app/ontology.project/project.module.ts":
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
const angular_split_1 = __webpack_require__("./node_modules/angular-split/esm2015/angular-split.js");
const project_component_1 = __webpack_require__("./src/app/ontology.project/project.component.ts");
const project_routing_module_1 = __webpack_require__("./src/app/ontology.project/project-routing.module.ts");
const layout_module_1 = __webpack_require__("./src/app/ontology.layout/layout.module.ts");
const layout_detail_module_1 = __webpack_require__("./src/app/ontology.layout.detail/layout.detail.module.ts");
const conversation_module_1 = __webpack_require__("./src/app/ontology.conversation/conversation.module.ts");
const shared_module_1 = __webpack_require__("./src/app/common/shared.module.ts");
const unfinished_topic_module_1 = __webpack_require__("./src/app/ontology.unfinished.topic/unfinished.topic.module.ts");
let ProjectModule = class ProjectModule {
};
ProjectModule = __decorate([
    core_1.NgModule({
        imports: [
            angular_split_1.AngularSplitModule,
            project_routing_module_1.ProjectRoutingModule,
            layout_module_1.LayoutModule,
            layout_detail_module_1.LayoutDetailModule,
            conversation_module_1.ConversationModule,
            shared_module_1.SharedModule,
            unfinished_topic_module_1.UnfinishedTopicModule
        ],
        declarations: [project_component_1.ProjectComponent]
    })
], ProjectModule);
exports.ProjectModule = ProjectModule;


/***/ })

});
//# sourceMappingURL=project.module.chunk.js.map