Covidtracer
-----------

This repository contains a privacy preserving contact tracer design for
COVID-19, read the [PDF](pp-contact-tracer.pdf).

### Key features

  - No personal data is shared with the system at any time. No constant user-ids, no cellphone records, no gps data, no phone 
    numbers, no names. No constant identifiers at all.
  - Users cannot be tracked through their broadcasted IDs: They are single-use, not linkable, and allow exactly one 
    notification to be sent without client authorization. The notification significance can be suppressed by dummy traffic.
  - Good user2user privacy to protect members of an infection chain in addition to random public contacts.
  - Low data traffic requirements. 120kB, 340kB or 10MB depending on client's privacy preferences. Data does not depend on 
    count of infections in the population.
  - Based on anonymized messaging protocols that have been known and understood for a long time.
  - Protocol allows the use of cheap dedicated embedded devices/wearables to prevent exposure of mobile phone to security 
    risks.
  - Processing requirements are low enough to support older smartphones.
  - Very low local storage requirements (<10MB).
  - Contact tracing avalanche controlled by health authorities.
  - Meaningful statistical input collected by the users themselves.


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
