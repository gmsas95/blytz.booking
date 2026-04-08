export const formatMoneyFromMinor = (amountMinor: number, currencyCode: string = 'USD') => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currencyCode,
  }).format(amountMinor / 100);
};

export const subtractMinorAmounts = (leftMinor: number, rightMinor: number) => leftMinor - rightMinor;
