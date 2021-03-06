# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased][unreleased]
### New
### Changed

## [0.2.4] - 2015-11-10
### Fixed
- Fixed better handling of Time.IsZero()

## [0.2.3] - 2015-11-10
### Fixed
- Better handling of Time.IsZero()

## [0.2.2] - 2015-10-27
### Fixed
- ModelValidate would return a non-nil value when it intended to return a nil value for no errors

## [0.2.1] - 2015-10-26
### Changed
- Wrapped dbr.Tx in an interface and updated Session implementation to return our wrap

## [0.2.0] - 2015-10-25
### New
- IsSet() bool added to all fields and is part of the field.Field interface
- field.Names can now return the intersection of two sets, Intersect()
### Changed
- Removed `Valid` field from non-nullable fields
- Fixed a few Marshal/Unmarshal edge cases


## [0.1.1] - 2015-10-19
### New
- Support for `ModelFields` to recurse into Anonymous/Embedded models to return all field.Names

## [0.1.0] - 2015-09-28
### New
- Connections are created with the database name, full database.tablename selectors for queries
  can be generated with `ModelTableName`.
### Changed
- Wrapped `dbr.Connection` and `dbr.Session` in our own local structs to add some context.
- Validation code refactored to two interfaces.


## [0.0.5] - 2015-09-22
- Float64 and NullFloat64 support

## [0.0.4] - 2015-09-14
### New
- Better handling of Null fields

## [0.0.3] - 2015-09-04
### New
- `ModelDirtyFields` added to get field.Names that are dirty
- `ModelLoadMap` can assign map[string]interface{} to a Model

### Changed
- Update and Insert use new PrimaryKeyer data to trim/customize queries
- JSON Marshal/Unmarshal improvements
