import "@typespec/http";

using TypeSpec.Http;

@service(#{ title: "Booking API" })
namespace Booking;

model Health {
  status: "ok";
}

model DatabaseUnavailable {
  code: 14;
  message: "database unavailable";
  details: unknown[];
}

@get
@route("/api/health")
op checkHealth():
  | Health
  | {
      @statusCode statusCode: 503;
      @body body: DatabaseUnavailable;
    };
