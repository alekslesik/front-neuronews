(function ($) {
	"use strict"

	$(window).on('scroll', function () {
		// Fixed Nav
		var wScroll = $(this).scrollTop();
		wScroll > $('header').height() ? $('#nav-header').addClass('fixed') : $('#nav-header').removeClass('fixed');

		// Back to top appear
		wScroll > 740 ? $('#back-to-top').addClass('active') : $('#back-to-top').removeClass('active')
	});

	// Back to top
	$('#back-to-top').on("click", function () {
		$('body,html').animate({
			scrollTop: 0
		}, 500);
	});

	// Mobile Toggle Btn
	$('#nav-header .nav-collapse-btn').on('click', function () {
		$('#main-nav').toggleClass('nav-collapse');
	});

	// Search Toggle Btn
	$('#nav-header .search-collapse-btn').on('click', function () {
		$(this).toggleClass('active');
		$('.search-form').toggleClass('search-collapse');
	});

	// Owl Carousel
	$('#owl-carousel-1').owlCarousel({
		loop: true,
		margin: 0,
		dots: false,
		nav: true,
		navText: ['<i class="fa fa-angle-left"></i>', '<i class="fa fa-angle-right"></i>'],
		autoplay: true,
		responsive: {
			0: {
				items: 1
			},
			992: {
				items: 2
			},
		}
	});

	$('#owl-carousel-2').owlCarousel({
		loop: false,
		margin: 15,
		dots: false,
		nav: true,
		navContainer: '#nav-carousel-2',
		navText: ['<i class="fa fa-angle-left"></i>', '<i class="fa fa-angle-right"></i>'],
		autoplay: false,
		responsive: {
			0: {
				items: 1
			},
			768: {
				items: 2
			},
			992: {
				items: 3
			},
		}
	});

	$('#owl-carousel-3').owlCarousel({
		items: 1,
		loop: true,
		margin: 0,
		dots: false,
		nav: true,
		navText: ['<i class="fa fa-angle-left"></i>', '<i class="fa fa-angle-right"></i>'],
		autoplay: true,
	});

	$('#owl-carousel-4').owlCarousel({
		items: 1,
		loop: true,
		margin: 0,
		dots: true,
		nav: false,
		autoplay: true,
	});

})(jQuery);

$(document).ready(function () {
	// Получаем текущий путь URL
	var currentPath = window.location.pathname;

	// Перебираем все элементы меню
	$('#main-nav .main-nav.nav.navbar-nav li').each(function () {
		// Получаем ссылку (href) из элемента меню
		var menuItemLink = $(this).find('a');

		// Если ссылка существует
		if (menuItemLink.length > 0) {
			// Получаем путь (href) из ссылки
			var menuItemPath = menuItemLink.attr('href');

			// Проверяем, совпадает ли текущий путь с путем элемента меню
			if (currentPath === menuItemPath) {
				// Если совпадает, добавляем класс "active" элементу меню
				$(this).addClass('active');
			} else {
				// Иначе, удаляем класс "active" (если он был установлен)
				$(this).removeClass('active');
			}
		}
	});
});


window.onscroll = function() {fixDiv()};

var fixedDivs = document.getElementsByClassName("fixedDiv");

var originalOffsets = [];
for (var i = 0; i < fixedDivs.length; i++) {
  originalOffsets[i] = fixedDivs[i].offsetTop;
}

function fixDiv() {
  for (var i = 0; i < fixedDivs.length; i++) {
    if (window.pageYOffset >= originalOffsets[i]) {
      fixedDivs[i].classList.add("fixed");
    } else {
      fixedDivs[i].classList.remove("fixed");
    }
  }
}
