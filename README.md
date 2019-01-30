[![Build Status](http://ci.endian.io/api/badges/endiangroup/specstack/status.svg?branch=develop)](http://ci.endian.io/endiangroup/specstack)
[![Go Report Card](https://goreportcard.com/badge/github.com/endiangroup/specstack)](https://goreportcard.com/report/github.com/endiangroup/specstack)
[![Endian homepage](https://endian.io/img/badge.svg)](https://endian.io)

# SpecStack. Specification as code. 

SpecStack is a work in progress multitool for managing Cucumber specifications. It has three broad parts:

1. **Meta data storage**. This allows the user to store additional textural information against stories or scenarios. It will eventually be used for implementation conversations, progress tracking and and some other project management applications.
2. **Specification refactoring**. This will allow users to rewords the language of the specification in an intuitive and powerful way.
3. Smart **Linting** tools, that can catch all sorts of problems in the specification before they arise.

## Testing and linting

Run `make list test`.
