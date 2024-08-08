CREATE TABLE IF NOT EXISTS public.blockchain
(
    key     uuid not null
        constraint blockchain_pk
            primary key,
    previous    text not null,
    merkley     text not null,
    noce        text not null,
    data        text not null
);

COMMENT ON TABLE public.blockchain IS 'Блокчейн';
COMMENT ON COLUMN public.blockchain.hash is 'Хэш предыдущего блока';
COMMENT ON COLUMN public.blockchain.merkley is 'Корень Меркли';
COMMENT ON COLUMN public.blockchain.data is 'Данные';

