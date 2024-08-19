CREATE TABLE IF NOT EXISTS public.blockchain
(
    key uuid NOT NULL,
    hash text COLLATE pg_catalog."default" NOT NULL,
    merkley text COLLATE pg_catalog."default" NOT NULL,
    nonce text COLLATE pg_catalog."default" NOT NULL,
    data text COLLATE pg_catalog."default" NOT NULL,
    date timestamp without time zone NOT NULL DEFAULT now(),
    CONSTRAINT blockchain_pk PRIMARY KEY (key)
)


COMMENT ON TABLE public.blockchain IS 'Блокчейн';
COMMENT ON COLUMN public.blockchain.hash IS 'Хэш предыдущего блока';
COMMENT ON COLUMN public.blockchain.merkley IS 'Корень Меркли';
COMMENT ON COLUMN public.blockchain.data IS 'Данные';

