TODO
----

- validators
  - default validators
  - validate method interface
  - explicitly call validate
  - model validate, valdiateFields([]string), validateAll()
  - validator state interface/struct
  - async validators

- sparse
  - dirty model flag on model
  - method on model to indicate changed fields
  - track origional values and changed values, condolidate on insert/update
  - fields specifically requested/set vs not

- package level functions for common functionality
  - IsDirty
  - ValidateAll
  - ValidateFields

- pre-select filter on select fields to reduce wire loads