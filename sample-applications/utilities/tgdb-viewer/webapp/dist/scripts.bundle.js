/**
 * UX Pattern Library v1.1.1 (http://uxpl.tibco.com.com)
 * Copyright © 2018 TIBCO Software Inc.
 */

// $(document).ready(function () {
//     console.log("Inside dropdown");
//     if($(".tc-dropdown > .dropdown-menu").length){
//         $(".tc-dropdown > .dropdown-menu").mCustomScrollbar({
//             scrollbarPosition:"inside"
//         });
//     }
// });
/**
 * Created by hede on 10/13/2017.
 */
// $(document).ready(function () {
//     console.log("Inside dropdown new");
//     /* if($(".tc-dropdown > .dropdown-menu").length){
//      $(".tc-dropdown > .dropdown-menu").mCustomScrollbar({
//      scrollbarPosition:"inside"
//      });
//      }
//      */
//     $("body").click(function(ev){
//         console.log("clicked out");
//         if ($(".dropdown-selector").length) {
//             console.log("inside");
//             if (!ev.target.parentNode.matches('.dropdown-selector')) {
//                 /*  if (!($(ev).hasClass('.dropdown-selector'))) {*/
//                 console.log("target not match");
//                 var dropdowns = document.getElementsByClassName("dropdown-menu");
//                 var i;
//                 for (i = 0; i < dropdowns.length; i++) {
//                     var openDropdown = dropdowns[i];
//                     if (openDropdown.classList.contains('tc-dd-show')) {
//                         console.log("openDropdown if");
//                         openDropdown.classList.remove('tc-dd-show');
//                     }
//                 }
//             }
//         }
//     });
//
// });

// var addClickEventListener = function (classname, callBackFunction) {
//     for (var i = 0; i < classname.length; i++) {
//         classname[i].addEventListener('click', callBackFunction, false);
//     }
// }
// var bodyCBFunction = function (event) {



// function openDropDownIconMenu() {
//     console.log("Inside openDropDown" + document.querySelector(".dropdown-icon-list"));
//     document.querySelector(".dropdown-icon-list").style.display = "block";
// }
// $(".tc-global-header-nav").click(function (event) {
//     event.stopPropagation();
//     $('.tc-container-background,.tc-panel-background').show();
//     $('.tc-icon-tray').addClass('active');
//     $('.tc-panel-container').addClass('slide-in');
// });
// if ($(".tc-global-header-container").length) {
//     $('body').click(function () {
//         $('.tc-container-background,.tc-panel-background').hide();
//         $('.tc-icon-tray').removeClass('active');
//         $('.tc-panel-container').removeClass('slide-in');
//         $(".tc-header-nav-item").removeClass('active');
//     });
// }
var getClassName = function (className) {
    return document.getElementsByClassName(className);
}
var tcHeaderNavItem = getClassName("tc-header-nav-item");
var tcGlobalHeaderNav = getClassName("tc-global-header-nav");

var navItemCbFn = function () {
    removeClass(tcHeaderNavItem, "active");
    var li = event.target.tagName.toLowerCase() === 'li' ? event.target : event.target.closest('li');
    li.classList.add("active");
    console.log("li", li);
    var item = li.classList[1];
    console.log("item", item);
    // document.getElementsByClassName("tc-panel-content " + item)[0].classList.add("active");
}
var addClass = function (element, className) {
    for (var i = 0; i < element.length; i++) {
        element[i].classList.add(className);
    }
}
var removeClass = function (element, className) {
    for (var i = 0; i < element.length; i++) {
        element[i].classList.remove(className);
    }
}


var navCBFunction = function (event) {
    event.stopPropagation();
    var showElement = document.getElementsByClassName('tc-container-background');
    var showElem_bg = document.getElementsByClassName('tc-panel-background');

    // for (var i = 0; i < showElement.length; i++) {
    showElement[0].style.display = 'block';
    // }
    // for (var j = 0; j < showElem_bg.length; i++) {
    showElem_bg[0].style.display = 'block';
    // }
    addClass(document.getElementsByClassName('tc-icon-tray'), "active")
    addClass(document.getElementsByClassName('tc-panel-container'), "slide-in")

}

var addClickEventListener = function (classname, callBackFunction) {
    for (var i = 0; i < classname.length; i++) {
        classname[i].addEventListener('click', callBackFunction, false);
    }
}
addClickEventListener(tcHeaderNavItem, navItemCbFn);
addClickEventListener(tcGlobalHeaderNav, navCBFunction);

/**
 * Created by hede on 11/23/2017.
 */
// $(document).ready(function () {
//     console.log("Inside modals");
//     //debugger;
//     $(".modal-large-body").mCustomScrollbar({
//         scrollbarPosition: "inside"
//     });
// });
// /**
//  * Created by hede on 9/28/2017.
//  */
//
//     $(document).ready(function () {
//         console.log("Inside notification");
//     $(".tc-details-sent").mCustomScrollbar({
//         scrollbarPosition:"outside"
//     });
//
// });
//

/**
 * Created by hede on 9/29/2017.
 */

function matchCountry(element, input) {
    var data = eval(element.getAttribute("data-var"));
    var reg = new RegExp(input.split('').join('\\w*').replace(/\W/, ""), 'i');
    if (data) {
        return data.filter(function (country) {
            if (country.value.match(reg)) {
                return country;
            }
        });
    } else {
        return '';
    }
}

function changeInput(element, val) {
    var autoCompleteResult = matchCountry(element, val);
    var temp = '';
    autoCompleteResult.forEach(function (entry) {
        temp += "<li onclick='selectItem(this)'>" + entry.value + "</li>";
    });
    document.getElementById("result").innerHTML = temp;
    document.getElementById("result").style.display = 'block';
}

function selectItem(item) {
    document.querySelector('.tc-search .tc-search-container > input').value = item.textContent;
    document.getElementById("result").style.display = 'none';
}


/**
 * Created by scheripa on 2/6/2018.
 */
/**
 * Created by hede on 10/13/2017.
 */

// $(document).ready(function () {
//     console.log("Inside dropdown new");
//     /* if($(".tc-dropdown > .dropdown-menu").length){
//      $(".tc-dropdown > .dropdown-menu").mCustomScrollbar({
//      scrollbarPosition:"inside"
//      });
//      }
//      */
//     $("body").click(function(ev){
//         console.log("clicked out");
//         if ($(".dropdown-selector").length) {
//             console.log("inside");
//             if (!ev.target.parentNode.matches('.dropdown-selector')) {
//                 /*  if (!($(ev).hasClass('.dropdown-selector'))) {*/
//                 console.log("target not match");
//                 var dropdowns = document.getElementsByClassName("dropdown-menu");
//                 var i;
//                 for (i = 0; i < dropdowns.length; i++) {
//                     var openDropdown = dropdowns[i];
//                     if (openDropdown.classList.contains('tc-dd-show')) {
//                         console.log("openDropdown if");
//                         openDropdown.classList.remove('tc-dd-show');
//                     }
//                 }
//             }
//         }
//     });
//
// });

document.querySelector('body').onclick = function (event) {
    if (event.target != undefined && !event.target.parentNode.matches('.tc-dropdown')
        && !event.target.parentNode.matches('.tc-buttons-dropdown')
        && !event.target.parentNode.matches('.dropdown-selector')
        && !event.target.parentNode.matches('.tc-search')
        && !event.target.parentNode.matches('.tc-search-container')) {
        showHideAllElements('.dropdown-menu, .tc-search-results', "none");
    }
};

// addClickEventListener(document.getElementsByTagName('body'),  );
function openDropDown(elem) {
    var dropDown = elem.nextElementSibling;
    if(dropDown.classList.contains('dropdown-menu')) {
        dropDown.style.display = "block";
    }
}

function showHideAllElements(element, displayType) {
    var elements = document.querySelectorAll(element);
    for (var index = 0; index < elements.length; index++) {
        elements[index].style.display = displayType;
    }
}




;
//# sourceMappingURL=scripts.bundle.js.map