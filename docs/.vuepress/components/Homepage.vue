<template lang="pug">
  div
    .search__container
      .search(@click="$emit('search', true)")
        .search__icon
          icon-search
        .search__text Search
    .h1 {{$frontmatter.title}}
    .intro
      .p {{$frontmatter.description}}
    .h2 Getting Started
    .p__alt Read all about the SDK or dive straight into the code with a tutorial.
    .features
      router-link(to="./intro/overview.html").features__item.features__item__light
        .features__item__image
          icon-spaceship.features__item__image__img
        .features__item__text
          .features__item__text__h2 Read
          .features__item__text__h1 Quick start
          .features__item__text__p A practical introduction to building blockchains with Ethermint.
          //- .features__item__text__tag
          //-   .features__item__text__tag__text 5 min
      a(href="https://tutorials.cosmos.network" target="_blank" rel="noopener").features__item.features__item__dark
        icon-arrow.features__item__icon
        .features__item__image
          icon-window.features__item__image__img
        .features__item__text
          .features__item__text__h2 Use
          .features__item__text__h1 Guides
          .features__item__text__p Follow guides to using polular Ethereum tools with Ethermint.
          //- .features__item__text__label View Tutorials
    .sections__wrapper
      .h2 Explore Ethermint
      .p__alt Get familiar with Ethermint and explore its main concepts.
      .sections
        router-link.sections__item(tag="a" :to="section.url" v-for="section in $frontmatter.sections")
          component(:is="`tm-icon-${section.icon}`").sections__item__icon
          .sections__item__wrapper
            .sections__item__title {{section.title}}
            .sections__item__desc {{section.desc}}
    .h2 Explore the Stack
    .p__alt Check out the docs for the various parts of the Ethermint stack.
    .stack
      a.stack__item(:href="item.url" v-for="item in $frontmatter.stack" :style="{'--accent-color': item.color, '--opacity': '5%'}")
        .stack__item__wrapper
          svg(width="17" height="16" viewBox="0 0 17 16" fill="none" xmlns="http://www.w3.org/2000/svg").stack__item__icon
            path(d="M1.07239 14.4697C0.779499 14.7626 0.779499 15.2374 1.07239 15.5303C1.36529 15.8232 1.84016 15.8232 2.13305 15.5303L1.07239 14.4697ZM15.7088 1.95457C16.0017 1.66168 16.0017 1.18681 15.7088 0.893912C15.4159 0.601019 14.941 0.601019 14.6482 0.893912L15.7088 1.95457ZM15.6027 1H16.3527C16.3527 0.585786 16.0169 0.25 15.6027 0.25V1ZM5.4209 0.25C5.00669 0.25 4.6709 0.585786 4.6709 1C4.6709 1.41421 5.00669 1.75 5.4209 1.75V0.25ZM14.8527 11.1818C14.8527 11.596 15.1885 11.9318 15.6027 11.9318C16.0169 11.9318 16.3527 11.596 16.3527 11.1818H14.8527ZM2.13305 15.5303L15.7088 1.95457L14.6482 0.893912L1.07239 14.4697L2.13305 15.5303ZM15.6027 0.25H5.4209V1.75H15.6027V0.25ZM16.3527 11.1818V1H14.8527V11.1818H16.3527Z" fill="#DADCE6")
          div
            .stack__item__h1 {{item.title}}
            .stack__item__p {{item.desc}}
    tm-help-support
</template>

<style lang="stylus" scoped>
/deep/
  .container h1
    margin-bottom 1.5rem

.search
  display flex
  align-items center
  color rgba(22, 25, 49, 0.65)
  padding-top 0.5rem
  width calc(var(--aside-width) - 6rem)
  cursor pointer

  &__container
    display flex
    justify-content flex-end
    margin-top 1rem
    margin-bottom 1rem

  &__icon
    width 1.5rem
    height 1.5rem
    fill #aaa
    margin-right 0.5rem

.intro
  width 100%
  max-width 800px

.h1
  font-size 2.5rem
  font-weight 600
  margin-bottom 1.5rem
  padding-top 2.5rem

.h2
  font-size 2rem
  font-weight 600
  margin-top 4.5rem
  margin-bottom 1rem

.p
  font-size 1.5rem
  line-height 36px

  &__alt
    margin-top 0.75rem
    margin-bottom 2rem
    font-size 1.25rem
    line-height 28px

.features
  display grid
  grid-template-columns repeat(auto-fit, minmax(300px, 1fr))
  gap 1.5rem
  margin-bottom 3rem
  margin-top 2.25rem

  &__item
    cursor pointer
    display grid
    grid-auto-flow column
    grid-template-columns 1fr 2fr
    box-shadow 0px 2px 4px rgba(22, 25, 49, 0.05), 0px 0px 1px rgba(22, 25, 49, 0.2), 0px 0.5px 0px rgba(22, 25, 49, 0.05)
    position relative
    border-radius .5rem
    background linear-gradient(302.07deg, #FFFFFF 48.96%, #EBEDFF 100%)
    transition box-shadow 0.25s
    
    &:hover
      box-shadow 0px 12px 24px rgba(22, 25, 49, 0.07), 0px 4px 8px rgba(22, 25, 49, 0.05), 0px 1px 0px rgba(22, 25, 49, 0.05)

    &__dark
      background linear-gradient(112.22deg, #161831 0%, #2E3148 100%)

    &__dark &__text__h2
      color white
      opacity .5

    &__dark &__text__h1
      color white

    &__dark &__text__p
      color white
      opacity .8

    &__dark &__text__label
      color white

    &__icon
      position absolute
      top 0
      right 0
      padding .75rem
      width 1rem
      height 1rem
      fill white
      opacity .35

    &__image
      display flex
      align-items center
      justify-content center

      &__img
        max-height 232px
        max-width 160px

    &__text
      padding 1.75rem 2rem 2rem
      display flex
      flex-direction column

      &__h2
        font-size .75rem
        letter-spacing 0.2em
        text-transform uppercase
        color #161931
        opacity .45
        margin-bottom .25rem

      &__h1
        font-size 1.25rem
        color #161931
        line-height 28px
        letter-spacing -0.01em
        font-weight 500
        margin-top .25rem
        margin-bottom .75rem

      &__p
        color rgba(22, 25, 49, 0.65)
        font-size .875rem
        letter-spacing 0.01em
        line-height 20px
        margin-bottom 1.5rem

      &__tag
        flex-grow 1
        display flex
        align-items flex-end

        &__text
          display inline-block
          background rgba(176, 180, 207, 0.2)
          border-radius .25rem
          padding-left .5rem
          padding-right .5rem
          color rgba(22, 25, 49, 0.65)
          letter-spacing 0.01em
          font-size .8125rem
          line-height 18px

      &__label
        text-transform uppercase
        font-weight 500
        letter-spacing 0.02em
        flex-grow 1
        display flex
        align-items flex-end

.sections
  display grid
  margin-top 3rem
  margin-bottom 4rem
  grid-template-columns repeat(auto-fit, minmax(300px, 1fr))
  gap 1.5rem

  &__item
    padding-left 4rem
    position relative
    color initial
    border-radius 0.5rem
    padding 1.5rem 1.5rem 1.5rem 4.5rem
    box-shadow 0px 2px 4px rgba(22, 25, 49, 0.05), 0px 0px 1px rgba(22, 25, 49, 0.2), 0px 0.5px 0px rgba(22, 25, 49, 0.05)
    transition box-shadow 0.25s

    &:hover
      box-shadow 0px 12px 24px rgba(22, 25, 49, 0.07), 0px 4px 8px rgba(22, 25, 49, 0.05), 0px 1px 0px rgba(22, 25, 49, 0.05)

    &__icon
      position absolute
      top 1rem
      left 1rem
      font-size 1.5rem
      display flex
      align-items center
      justify-content center
      width 2.5rem
      height 2.5rem

    &__title
      font-weight 500
      margin-bottom 0.5rem

    &__desc
      font-size 0.875rem
      color rgba(22, 25, 49, 0.65)

.stack
  display grid
  gap 1.5rem
  grid-template-columns repeat(auto-fit, minmax(300px, 1fr))
  margin-bottom 3rem

  &__item
    position relative
    min-height 120px
    display flex
    align-items center
    padding 2rem 1rem
    border-radius 0.5rem
    box-shadow inset 0 0 0 1px rgba(176, 180, 207, 0.2)
    color #161931
    background white
    transition box-shadow 0.25s

    &:hover
      box-shadow inset 0 0 0 1px rgba(176, 180, 207, 0.2), 0px 12px 24px rgba(22, 25, 49, 0.07), 0px 4px 8px rgba(22, 25, 49, 0.05), 0px 1px 0px rgba(22, 25, 49, 0.05)

    &__icon
      position absolute
      top 0
      right 0
      padding 1rem

    &__h1
      font-size 1.25rem
      margin-bottom 0.75rem
      font-weight 500

    &__p
      font-size 0.875rem
      color rgba(22, 25, 49, 0.65)
      line-height 1.25rem

    &__wrapper
      display grid
      grid-auto-flow column
      grid-template-columns 30% 1fr
      gap 1rem

    &:before
      position absolute
      top 0
      left 0
      content ''
      width 50%
      height 100%
      background linear-gradient(to right, var(--accent-color), rgba(255, 255, 255, 0))
      border-radius 0.5rem
      opacity 0.1

    &__logo
      height 72px
      width auto

@media screen and (max-width: 1136px)
  .p
    font-size 1.25rem
    line-height 28px

@media screen and (max-width: 832px)
  .h1
    padding-top 3.5rem

  .search__container
    display none

@media screen and (max-width: 752px)
  .search
    display none

@media screen and (max-width: 500px)
  .sections
    gap 0.5rem
    margin-bottom 0
    margin-top 2rem
    grid-template-columns repeat(auto-fit, minmax(220px, 1fr))

    &__item
      margin-bottom 0
      padding-left 4rem
      box-shadow none
      padding-bottom 0

      &__icon
        top 1rem
        left 0.75rem

      &__wrapper
        padding-bottom 1.5rem
        border-bottom 1px solid rgba(140, 145, 177, 0.32)

      &:last-child .sections__item__wrapper
        border-bottom none

    &__wrapper
      position relative
      padding 0.1px 1rem
      background white
      border-radius 0.5rem

      &:after
        position absolute
        content ''
        top 0
        left 0
        right 0
        bottom 0
        border-radius 0.5rem
        box-shadow 0px 24px 40px rgba(22, 25, 49, 0.1), 0px 10px 16px rgba(22, 25, 49, 0.08), 0px 1px 0px rgba(22, 25, 49, 0.05)

  .stack
    gap 0.5rem
    grid-template-columns repeat(auto-fill, minmax(240px, 1fr))

    &__item
      &:after
        position absolute
        content ''
        top 0
        left 0
        right 0
        bottom 0
        border-radius 0.5rem
        box-shadow 0px 24px 40px rgba(22, 25, 49, 0.1), 0px 10px 16px rgba(22, 25, 49, 0.08), 0px 1px 0px rgba(22, 25, 49, 0.05)

@media screen and (max-width: 480px)
  .h1
    font-size 2rem

  .h2
    font-size 1.5rem

  .p__alt
    font-size 1rem
    line-height 1.5rem

  .features
    grid-template-columns repeat(auto-fit, minmax(240px, 1fr))

    &__item
      display block

      &__image
        padding-top 1rem

      &__text
        padding 1.5rem
</style>
