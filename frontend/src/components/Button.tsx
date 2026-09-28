import type { ButtonHTMLAttributes } from 'react';

type Variant = 'primary' | 'secondary' | 'danger';

const variantClasses: Record<Variant, string> = {
  primary: 'bg-primary text-white hover:opacity-90 disabled:bg-muted',
  secondary:
    'bg-surface text-ink border border-control hover:bg-background disabled:text-muted',
  danger: 'bg-danger text-white hover:opacity-90 disabled:bg-muted',
};

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
}

/** Tombol platform: radius 8 px, tinggi min 44 px (DESIGN §4.2). */
export function Button({ variant = 'primary', className = '', type, ...rest }: ButtonProps) {
  return (
    <button
      type={type ?? 'button'}
      className={`inline-flex min-h-11 items-center justify-center gap-2 rounded-control px-4 text-label font-semibold transition-opacity duration-150 disabled:cursor-not-allowed disabled:opacity-60 ${variantClasses[variant]} ${className}`}
      {...rest}
    />
  );
}
