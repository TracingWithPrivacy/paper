Covidtracer
-----------

This repository contains a privacy preserving contact tracer design for
COVID-19, read the [PDF](pp-contact-tracer.pdf).

### Dependencies

-   make
-   pandoc
-   pandoc-citeproc
-   pdflatex

Most dependencies can be easily installed using Apt (Linux), Homebrew
(macOS), or Chocalety (Windows).

Apt (Linux) one-liner:

    sudo apt install make pandoc pandoc-citeproc texlive

Homebrew (macOS) one-liner:

    brew install make pandoc pandoc-citeproc && brew cask install basictex

### Build

    make
