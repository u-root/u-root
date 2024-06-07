import { gsap } from 'gsap';
import ScrollTrigger from 'gsap/ScrollTrigger';
gsap.registerPlugin(ScrollTrigger);

const mm = gsap.matchMedia(),
  bp = 600;

const homeHero = document.querySelector('.home-hero');

const keyVisualWrapper = document.querySelector('.home-hero__key-visual');
gsap.set(keyVisualWrapper, { height: 0, autoAlpha: 0 });

// Calculate relative offset to the top of the page in percentages
const startOffsetPercent = () => {
  const heroMarginTop = homeHero.getBoundingClientRect().top;
  return Math.ceil((heroMarginTop / homeHero.offsetHeight) * 100);
};

/*
 * Move, fade and scale key visual on scroll
 */
gsap.to(keyVisualWrapper, {
  height: 'auto',
  autoAlpha: 1,
  duration: 1,
  ease: 'linear',
  scrollTrigger: {
    trigger: homeHero,
    pin: true,
    start: () => `-=${startOffsetPercent()}%`,
    end: '+=45%',
    scrub: 0.25,
    invalidateOnRefresh: true,
    //markers: true,
  },
});

/*
 * Move and fade box items on scroll
 */
const keyVisualItems = document.querySelectorAll('.key-visual__item');
gsap.set(keyVisualItems, { y: 0, autoAlpha: 1 });

const homeHeroHeadline = document.querySelector('.home-hero__headline');

for (const [i, item] of keyVisualItems.entries()) {
  // Stagger animation for keyVisualItems
  const offset = 100; // Adjust this offset value as needed
  const distance = 70; // Adjust this distance value as needed

  let start = `+=${i * distance + offset}%`;
  console.log('start', start);

  gsap.to(item, {
    scrollTrigger: {
      trigger: homeHeroHeadline,
      start: start,
      end: '+=10%',
      scrub: 0.25,
      invalidateOnRefresh: true,
      //markers: true,
    },
    duration: 1,
    y: -60,
    autoAlpha: 0,
    ease: 'linear',
    id: `keyVisualItem-${i + 1}`,
  });
}

/*
 * Move and fade feature illustrations on scroll
 */

const featureSection = document.querySelector('.features');
const featuresCardIllustrations = featureSection.querySelectorAll(
  '.features__card svg'
);
gsap.set(featuresCardIllustrations, {
  y: -75,
  autoAlpha: 0,
});

const featuresCards = document.querySelectorAll('.features__card');
//const cardsReverseOrder = Array.from(featuresCards).reverse();

mm.add(
  {
    // set up any number of arbitrarily-named conditions. The function below will be called when ANY of them match.
    isDesktop: `(min-width: ${bp}px)`,
    isMobile: `(max-width: ${bp - 1}px)`,
  },
  (context) => {
    // context.conditions has a boolean property for each condition defined above indicating if it's matched or not.
    let { isDesktop, isMobile } = context.conditions;

    for (const [i, item] of featuresCards.entries()) {
      // Stagger animation for feature illustrations
      const offset = 35; // Adjust this offset value as needed
      const distance = 10; // Adjust this distance value as needed

      let start = `+=${(i + 1) * distance + offset}%`;
      //let end = `+=${(i + 1) * distance + offset + 100}%`;

      let illustration = item.querySelector('svg');

      gsap.to(illustration, {
        autoAlpha: 1,
        y: 0,
        duration: 1,
        ease: 'linear',
        id: `feature-${i + 1}`,
        scrollTrigger: {
          trigger: isDesktop ? homeHero : item,
          start: isDesktop ? start : '-=200%',
          end: isDesktop ? '+=10%' : '-=150%',
          scrub: 0.25,
          invalidateOnRefresh: true,
          //markers: true,
        },
      });
    }
  }
);
