import { useId } from 'react';
import type { InputHTMLAttributes } from 'react';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  helperText?: string;
  errorText?: string;
}

/** Input berlabel: batas kontrol #64748B, radius 8 px, tinggi min 44 px (DESIGN §4). */
export function Input({ label, helperText, errorText, id, ...rest }: InputProps) {
  const autoId = useId();
  const inputId = id ?? autoId;
  const describedBy = errorText ? `${inputId}-error` : helperText ? `${inputId}-helper` : undefined;

  return (
    <div className="flex flex-col gap-1">
      <label htmlFor={inputId} className="text-label font-semibold text-ink">
        {label}
      </label>
      <input
        id={inputId}
        aria-invalid={errorText ? true : undefined}
        aria-describedby={describedBy}
        className={`min-h-11 rounded-control border bg-surface px-3 text-body text-ink placeholder:text-muted ${
          errorText ? 'border-danger' : 'border-control'
        }`}
        {...rest}
      />
      {helperText && !errorText && (
        <p id={`${inputId}-helper`} className="text-label text-muted">
          {helperText}
        </p>
      )}
      {errorText && (
        <p id={`${inputId}-error`} className="text-label text-danger" role="alert">
          {errorText}
        </p>
      )}
    </div>
  );
}
