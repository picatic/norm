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

Examples
========

Model
-----
```golang
type Post struct {
  Id field.Int64
  Title field.String
  Content field.String
  AuthorId field.Int64
  Created field.Time
}
```

Select a Model
--------------

```golang
var post Post = Post{}
post.Id.Scan(4384)

err := norm.NewSelect(dbrSession, post, nil).Where("id = ?", post.Id.Int64).LoadStruct(&post) // select all defined fields

err = norm.NewSelect(dbrSession, post, field.FieldNames{"id", "title"}).Where("id = ?", post.Id.Int64).LoadStruct(&post) // only those fields
```

Select Models
-------------

```golang
var posts []Post = make([]Post, 0)

err := norm.NewSelect(dbrSession, post, nil).Where("created > '2015-01-01").LoadStructs(&posts)
```

Insert Model
------------

```golang
var post Post = Post{}
post.Title.Scan("First Post")
post.Context.Scan("...")
post.AuthorId.Scan(1)
post.Created.Scan(time.Time.Now())

result, err := norm.NewInsert(dbrSession, post, nil).Exec()

fmt.Printf("InsertId: %s", result.Value())
```

Update Model
------------

```golang

// continued from Insert
post.Content.Scan("modified content")

result, err := norm.NewUpdate(dbrSession, post, nil).Exec()

```

FAQ
===

Q: Why did you wrap all types?
A: That was the only way to get dirty model detection without generating Getter/Setters. It also happens to allow us to pretend to do sparse models.

Q: Why do I have to `Scan` values when I could set their value directy. e.g. model.Id.Int64 = 1 vs model.Id.Scan(1)
A: Scan currently contains logic to mark the shadow value if not already set. This is key in determining if we are working with dirty fields and thus dirty models.

Q: Why do I have to set the LastInsertId of an UPDATE to my model manually?
A: NORM does not persume you care of want that functionality. A simple utility function can wrap that logic for you to meet your needs.

```golang
// Insert and update id on model
func InsertAndUpdateId(sess dbr.ConnSession, model Model, fields field.FieldNames) error {
  result, err := norm.NewInsert(sess, model, fields).Exec()
  if err != nil {
    return err
  }
  lastId, err := resp.LastInsertId()
  if err != nil {
   return err
  }
  field := norm.ModelGetField(model.PrimaryKeyFieldName())
  field.Scan(lastId)
  return nil
}
```

Q: Have you considered implementing before and after hooks for models?
A: Yes, but given how we use dbr and imagined NORM as a lazy ORM it is not something will will build in the ORM. We may provide an example that suggests a pattern of interfaces to do hooks, but ultimately you would need to call those hooks by convention.

Q: I want feature XYZ?
A: Post an issue as a proposal. Outline what the feature is, how it would work and be implemented and how it fits into NORM. We will discuss and go from there.

Q: I found a bug!
A: Not a question, but submit an issue with the details. Ideally a test that triggers the bug. Even more ideally, a test and patch for the bug.

Q: Why did you use dbr?
A: Flexability of programatically building queries or just making raw queries.