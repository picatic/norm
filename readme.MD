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