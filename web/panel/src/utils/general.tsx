import { notifications } from "@mantine/notifications";
import { BadgeCheck, CircleAlert, CircleX, Info } from "lucide-react";

// Type utility as type guards and type transformations function.
export function isFetchSuccess<T>(response: {
  data: T;
  code: number;
  error?: undefined;
} | {
  error: any;
  data?: undefined;
  code?: undefined;
}): response is { data: T; code: number; error: undefined } {
  return !response.error;
}

type ToastShow = {
  title: string;
  message: string;
  loading?: boolean;
  autoClose?: number;
  withCloseButton?: boolean;
  status?: 'ok' | 'issue' | 'warn' | 'info';
};

export function showToast({ title, message, autoClose, status, loading, withCloseButton }: ToastShow) {
  autoClose = autoClose || 4000;
  const icon = status === 'issue' ? <CircleX size={24} />
    : status === 'warn' ? <CircleAlert />
    : status === "info" ? <Info />
    : <BadgeCheck />;

  notifications.show({
    color: status === 'issue' ? 'red' : status === 'warn' ? 'yellow' : status === 'info' ? 'blue' : 'green',
    withCloseButton,
    mod: { status },
    autoClose,
    loading,
    message,
    title,
    icon,
  });
}
