define(function () {
	"use strict";
	function ariSkyParallax() {
		$('body').parallax({
			'elements': [
				{
					'selector': '.bg-1',
					'properties': {
						'x': {
							'background-position-x': {
								'initial': 0,
								'multiplier': 0.02,
								'invert': true
							}
						}
					}
				},
				{
					'selector': '.bg-2',
					'properties': {
						'x': {
							'background-position-x': {
								'initial': 0,
								'multiplier': 0.06,
								'invert': true
							}
						}
					}
				},
				{
					'selector': '.bg-3',
					'properties': {
						'x': {
							'background-position-x': {
								'initial': 0,
								'multiplier': 0.2,
								'invert': true
							}
						}
					}
				}
			]
		});

		$(window).scroll(function () {

			var a = 50;
			var pos = $(window).scrollTop();
			if ($(window).width()>=1200){
				if (pos>a){
					$("#ariGirl").removeClass("animated bounceIn");
					$("#ariGirl").addClass("animated fadeOut");
				}
				else{
					$("#ariGirl").removeClass("animated fadeOut");
					$("#ariGirl").addClass("animated bounceIn");
				}
			}
			else{
				$(".rightFiesta").css({
					display: 'none'
				});
			}
			$(window).resize(function(){
				if ($(window).width()<1200) {
					$(".rightFiesta").css({
						display: 'none'
					});
				}
				else{
					$(".rightFiesta").css({
						display: 'block'
					});
				}
			});
		});
	}

	return $(function () {
		ariSkyParallax();
	});
});