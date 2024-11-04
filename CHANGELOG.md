# Changelog

## [v0.13.0](https://github.com/k1LoW/expand/compare/v0.12.0...v0.13.0) - 2024-11-04
### Other Changes
- chore(deps): bump the dependencies group with 4 updates by @dependabot in https://github.com/k1LoW/expand/pull/48
- chore(deps): bump golangci/golangci-lint-action from 2 to 6 in the dependencies group by @dependabot in https://github.com/k1LoW/expand/pull/46

## [v0.12.0](https://github.com/k1LoW/expand/compare/v0.11.0...v0.12.0) - 2024-03-08
### Breaking Changes 🛠
- Escape/unescape delimiters by @k1LoW in https://github.com/k1LoW/expand/pull/44

## [v0.11.0](https://github.com/k1LoW/expand/compare/v0.10.1...v0.11.0) - 2023-12-30
### Other Changes
- Update expr (change org) by @k1LoW in https://github.com/k1LoW/expand/pull/43

## [v0.10.1](https://github.com/k1LoW/expand/compare/v0.10.0...v0.10.1) - 2023-10-30

## [v0.10.0](https://github.com/k1LoW/expand/compare/v0.9.0...v0.10.0) - 2023-10-30
### Breaking Changes 🛠
- String types "true" and "false" are converted as string types as much as possible. by @k1LoW in https://github.com/k1LoW/expand/pull/39

## [v0.9.0](https://github.com/k1LoW/expand/compare/v0.8.0...v0.9.0) - 2023-10-09
### Breaking Changes 🛠
- Update goccy/go-yaml to v0.11.2 by @k1LoW in https://github.com/k1LoW/expand/pull/35
### Other Changes
- Fix ExprRepFn to stringify correctly. by @k1LoW in https://github.com/k1LoW/expand/pull/37
- Update github.com/antonmedv/expr to v1.15.3 by @k1LoW in https://github.com/k1LoW/expand/pull/38

## [v0.8.0](https://github.com/k1LoW/expand/compare/v0.7.0...v0.8.0) - 2023-05-07
- Refactor code by @k1LoW in https://github.com/k1LoW/expand/pull/30
- [BREAKING CHANGE] If the expanded result is a Map or Slice, it is not quoted to be interpreted as inline YAML. by @k1LoW in https://github.com/k1LoW/expand/pull/32
- [BEAKING CHANGE] Fix sig of ReplaceYAML and Add options ( ReplaceMapKey, QuoteCollection ) by @k1LoW in https://github.com/k1LoW/expand/pull/33
- Update pkgs by @k1LoW in https://github.com/k1LoW/expand/pull/34

## [v0.7.0](https://github.com/k1LoW/expand/compare/v0.6.1...v0.7.0) - 2023-03-04
- Avoid duplicate quotes heuristically by @k1LoW in https://github.com/k1LoW/expand/pull/27
- Bump golang.org/x/sys from 0.0.0-20220406163625-3f8b81556e12 to 0.1.0 by @dependabot in https://github.com/k1LoW/expand/pull/29

## [v0.6.1](https://github.com/k1LoW/expand/compare/v0.6.0...v0.6.1) - 2023-03-04
- Fix oversight from https://github.com/k1LoW/expand/pull/24 by @k1LoW in https://github.com/k1LoW/expand/pull/25

## [v0.6.0](https://github.com/k1LoW/expand/compare/v0.5.6...v0.6.0) - 2023-03-03
- Stop giving unnecessary trailing newlines by @k1LoW in https://github.com/k1LoW/expand/pull/22
- If there is a line break in the result of the conversion of what was one line, quote it. by @k1LoW in https://github.com/k1LoW/expand/pull/24

## [v0.5.6](https://github.com/k1LoW/expand/compare/v0.5.5...v0.5.6) - 2023-02-04
- Bump up go and pkgs version by @k2tzumi in https://github.com/k1LoW/expand/pull/20

## [v0.5.5](https://github.com/k1LoW/expand/compare/v0.5.4...v0.5.5) - 2022-11-21
- Fix InterpolateRepFn ( handling env placeholder ) by @k1LoW in https://github.com/k1LoW/expand/pull/18

## [v0.5.4](https://github.com/k1LoW/expand/compare/v0.5.3...v0.5.4) - 2022-11-14
- Fix quoting strings by @k1LoW in https://github.com/k1LoW/expand/pull/16

## [v0.5.3](https://github.com/k1LoW/expand/compare/v0.5.2...v0.5.3) - 2022-10-05
- Fix quoting string in ReplaceYAML() by @k1LoW in https://github.com/k1LoW/expand/pull/14

## [v0.5.2](https://github.com/k1LoW/expand/compare/v0.5.1...v0.5.2) - 2022-10-05
- Support same delims by @k1LoW in https://github.com/k1LoW/expand/pull/11
- Fix trySubstr() by @k1LoW in https://github.com/k1LoW/expand/pull/13

## [v0.5.1](https://github.com/k1LoW/expand/compare/v0.5.0...v0.5.1) - 2022-10-04
- Fix ExprRepFn by @k1LoW in https://github.com/k1LoW/expand/pull/9

## [v0.5.0](https://github.com/k1LoW/expand/compare/v0.4.0...v0.5.0) - 2022-10-04
- [BREAKING] Introduce `type repFn func(in string) (string, error)` and Added error to the return values of ReplaceYAML() by @k1LoW in https://github.com/k1LoW/expand/pull/5
- Use tagpr by @k1LoW in https://github.com/k1LoW/expand/pull/7
- Add ExprRepFn by @k1LoW in https://github.com/k1LoW/expand/pull/6

## [v0.4.0](https://github.com/k1LoW/expand/compare/v0.3.0...v0.4.0) (2022-07-18)

* Support bash like interpolate [#4](https://github.com/k1LoW/expand/pull/4) ([k2tzumi](https://github.com/k2tzumi))

## [v0.3.0](https://github.com/k1LoW/expand/compare/v0.2.0...v0.3.0) (2022-03-06)

* Support key replace [#3](https://github.com/k1LoW/expand/pull/3) ([k1LoW](https://github.com/k1LoW))
* Add ReplaceYAML [#2](https://github.com/k1LoW/expand/pull/2) ([k1LoW](https://github.com/k1LoW))

## [v0.2.0](https://github.com/k1LoW/expand/compare/v0.1.0...v0.2.0) (2022-03-05)

* Update pkgs [#1](https://github.com/k1LoW/expand/pull/1) ([k1LoW](https://github.com/k1LoW))

## [v0.1.0](https://github.com/k1LoW/expand/compare/0c0882c8638e...v0.1.0) (2021-10-09)
