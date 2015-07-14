```
  _   _  ____  _____  __  __ 
 | \ | |/ __ \|  __ \|  \/  |
 |  \| | |  | | |__) | \  / |
 | . ` | |  | |  _  /| |\/| |
 | |\  | |__| | | \ \| |  | |
 |_| \_|\____/|_|  \_\_|  |_|
```

NORM is an ORM that tries not to be much like an ORM.

Features
========

* Forces you to make models that uses our fields.
  * Wraps fields on your models so you know if you changed them.
  * Save only the fields that changed.
  * Only Validate fields that changed.
* Validation interface that probably does all the things.
* Barely makes queries for you, actually that may be an overstatement.
* Built on dbr, so you can make your own queries that are better than anything we could make for you.

Why?
====

Because we wanted something like an ORM, but got out of our way and let us drive. Making queries is likely still best
done by humans (for now...). Most ORM's try to do too much, and eventually block you from doing the complex business 
logic you end up needing to do 6 months from now.

# ORM

## Resources

- https://github.com/lann/squirrel

## milestone targets
- Norm (Spike Interfaces)
- Data Validation (Spike Interfaces)
- Controller Input
- Controller Output 
- Authorization
- Implement Some Email Templating Endpoint
- Formalize Norm
- Generate types


## TODO
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

- generator
  - Set<FieldName>/<FieldName>

- package level functions for common functionality
  - IsDirty
  - ValidateAll
  - ValidateFields
  - etc

- Interfaces on Shadow Models
  - Valdiation default fields e.g. map[string][]Validators
  - Validation custom method

- types
  - DBR/sql NullWrappers

- Controllers
  - Convention named actions: <ModelPlural><Action>?

- permissions
  - filter output results (to be confirmed)

- pre-select filter on select fields to reduce wire loads

- auth
  - header: x-picatic-access-key
    - user specific access_key
  - query param: x-picatic-access-key
    - user specific access key
  - oauth2 by auth bearer


- input format
  - single model post/save/insert 
  - future feature: noonce

- output formats
  - object: "list" for multiple items
  - url: '/v2/users'
  - if set, is_more: bool
  - data -> [{...}...]
  - relations, do not return relations for now