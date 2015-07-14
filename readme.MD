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