this is a common framework for golang based servers and clients in guru.

The repository includes common code - but not APIs. (see proto repository).

examples are provided in the repository skel-go.

the tarfile includes golang.gurusys.co.uk/go-framework.

To add it to your module:

cd src/golang.gurusys.co.uk/vendor/golang.gurusys.co.uk && tar -jxvf /tmp/go-framework.tar.bz2

or, if you prefer to download, untar it without review:

cd src/golang.gurusys.co.uk/vendor/golang.gurusys.co.uk
wget -O - http://buildrepo.guru.localdomain/builds/go-framework/master/latest/dist/go-framework.tar.bz2 | tar -jx

you also need the apis:
wget -O - http://buildrepo.guru.localdomain/builds/proto/master/latest/dist/golang-api-stubs.tar.bz2 | tar -jx


# migrations from an earlier format:

There were earlier, incorrect, references.
This will lead to "cannot find package..."
Below are common mistakes and how to fix them:

github.com/GuruSystems/go-framework/server -> golang.gurusys.co.uk/go-framework/server

gitlab.gurusys.co.uk/reviewrequired/proto/ -> golang.gurusys.co.uk/apis/
