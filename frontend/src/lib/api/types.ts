// Tipe envelope API_SPEC §2 yang dipakai klien HTTP frontend.

export interface ApiMeta {
  request_id: string;
  page?: number;
  per_page?: number;
  total?: number;
}

export interface ApiEnvelope<T> {
  data: T;
  meta: ApiMeta;
}

export interface ApiErrorBody {
  code: string;
  message: string;
  fields?: Record<string, string>;
  current_version?: number;
}

export interface ApiErrorEnvelope {
  error: ApiErrorBody;
  meta: ApiMeta;
}
