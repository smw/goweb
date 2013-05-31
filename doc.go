// Goweb v2 - A simple, powerful web framework for Go.
//
// Overview
//
// Goweb follows the MVC pattern where you build controller objects, and map them to routes
// of URLs using the goweb.MapController function.  Controller objects should adhear to one or
// more of the controllers.Restful* interfaces in order to be mapped correctly.
//
// If you are not following RESTful patterns, you can do custom routing using the goweb.Map function.
//
// Example
//
// Your controllers, following a RESTful pattern, might look like this:
//
//     // PeopleController controls the 'people' resources.
//     type PeopleController struct {}
//
//     // ReadMany reads many people.
//     func (c *PeopleController) ReadMany(ctx context.Context) error {
//
//       // TODO: show all people
//
//     }
//
//     // Read reads one person.
//     func (c *PeopleController) Read(id string, ctx context.Context) error {
//
//       // TODO: show one person
//
//     }
//
//     // Create creates a new person.
//     func (c *PeopleController) Create(ctx context.Context) error {
//
//       // TODO: create a person, and redirect to the Read method
//
//     }
//
// In the above controller code, we are providing three RESTful methods, Read, ReadMany and Create.
//
// To map this in Goweb, we use the MapController function like this:
//
//     mapErr := goweb.MapController(&PeopleController{})
//
// This will map the two functions (since they follow the standards defined in the controllers package)
// to the appropriate RESTful URLs:
//
//     GET /people - ReadMany
//     GET /people/{id} - Read
//     POST /people - Create
//
// To add more RESTful features (like Update and Delete), you just have to add the relevant methods to the
// controller and Goweb will do the rest for you.
//
// For a full list of the supported methods, see the MapController function.
//
// Mapping paths
//
// If you want to map specific paths with Goweb, you can use the `goweb.Map` function with
// the following kinds of paths:
//
//
//     /literals - Normal text will be considered literal, meaning the path would
//                 have to match it exactly.
//
//     {variable} - Curly braces represent a segment of path that can be anything,
//                  and which will be made available via the `context.PathParam` method.
//
//     [optional-variable] - Square braces represent an optional segment.  If present,
//                           will be made available via the `context.PathParam` method
//                           but otherwise the path will still be a match.
//
//     * - A single `*` is similar to using `{}`, except the value is ignored.
//
//     *** - Three `*`'s matches anything in this segment, and any subsequent segments.
//           For example, `/people/***` would match `/people`, `/people/123` and `/people/123/books/456`.
//
// For some real examples of mapping paths, see the goweb.Map function, or check out the
// example_webapp in the code.
//
// Responding
//
// Goweb makes it easy to respond to requests using an extensible Responder pattern.
//
// If you're building an API, you can respond using the `goweb.API` object, if you just
// need to respond in a normal HTTP manner, you can use the `goweb.Respond` object.
//
// For details on how to make normal HTTP responses, see http://godoc.org/github.com/stretchrcom/goweb/responders#HTTPResponder
//
// For details on how to make API responses, see http://godoc.org/github.com/stretchrcom/goweb/responders#APIResponder
package goweb