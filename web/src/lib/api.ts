import createClient from "openapi-fetch";
import type { paths } from "@/gen/schema";

export const api = createClient<paths>({ baseUrl: "" });
