#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

#include <bliss_b_keys.h>
#include <bliss_b_signatures.h>
#include <sampler.h>

void test_verify(entropy_t *entropy, bliss_param_t *params) {
	int i;
	bliss_private_key_t private_key;
	bliss_public_key_t public_key;
	bliss_signature_t sig;
	const char *s = "Hello world";
	size_t msg_sz = strlen(s);
	const uint8_t *msg = (const uint8_t*)s;
	int32_t res;

	/* Generate private key */
	bliss_b_private_key_gen(&private_key, params->kind, entropy);
	bliss_b_public_key_extract(&public_key, &private_key);

	/* Output private key */
	printf("Private Key:\n");
	printf("s1: ");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.s1[i]);
	printf("\n");
	printf("s2: ");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.s2[i]);
	printf("\n");
	printf("a: ");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.a[i]);
	printf("\n");

	/* This is done for the sole purpose of allocating memory
	 * fo signature */
	bliss_b_sign(&sig, &private_key, msg, msg_sz, entropy);
	printf("z1: ");
	for(i = 0; i < params->n; i++)
		scanf("%d",&(sig.z1[i]));
	printf("z2: ");
	for(i = 0; i < params->n; i++)
		scanf("%d",&(sig.z2[i]));
	printf("c: ");
	for(i = 0; i < params->kappa; i++)
		scanf("%u",&(sig.c[i]));

	/* Verify */
	printf("Verify:\n");
	res = bliss_b_verify(&sig,&public_key,msg,msg_sz);
	printf("Verify res: %d\n",res);

	bliss_b_private_key_delete(&private_key);
	bliss_signature_delete(&sig);
}

void test_sign(entropy_t *entropy, bliss_param_t *params) {
	int i;

	bliss_private_key_t private_key;
	bliss_public_key_t public_key;
	bliss_signature_t sig;
	const char *s = "Hello world";
	size_t msg_sz = strlen(s);
	const uint8_t *msg = (const uint8_t*)s;
	int32_t res;

	/* Generate private key */
	bliss_b_private_key_gen(&private_key, params->kind, entropy);
	bliss_b_public_key_extract(&public_key, &private_key);

	/* Output private key */
	printf("Private Key:\n");
	printf("s1: ");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.s1[i]);
	printf("\n");
	printf("s2: ");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.s2[i]);
	printf("\n");
	printf("a: ");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.a[i]);
	printf("\n");

	/* Sign */
	printf("Sign:\n");
	bliss_b_sign(&sig, &private_key, msg, msg_sz, entropy);
	printf("Signature:\n");
	printf("z1: ");
	for(i = 0; i < params->n; i++)
		printf("%d ",sig.z1[i]);
	printf("\n");
	printf("z2: ");
	for(i = 0; i < params->n; i++)
		printf("%d ",sig.z2[i]);
	printf("\n");
	printf("c: ");
	for(i = 0; i < params->kappa; i++)
		printf("%d ",sig.c[i]);
	printf("\n");

	/* Verify */
	printf("Verify:\n");
	res = bliss_b_verify(&sig,&public_key,msg,msg_sz);
	printf("Verify res: %d\n",res);

	bliss_b_private_key_delete(&private_key);
	bliss_signature_delete(&sig);
}

void gen_private_key(entropy_t *entropy, bliss_param_t *params) {
	int i;

	bliss_private_key_t private_key;

	/* Generate private key */
	bliss_b_private_key_gen(&private_key, params->kind, entropy);

	/* Output private key */
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.s1[i]);
	printf("\n");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.s2[i]);
	printf("\n");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.a[i]);
	printf("\n");

	bliss_b_private_key_delete(&private_key);

}

void usage(char *argv[]) {
	fprintf(stderr,"Usage: %s subcommand [options]\n",argv[0]);
	fprintf(stderr,"  subcommands\n");
	fprintf(stderr,"    keygen\n");
	fprintf(stderr,"    sign\n");
	fprintf(stderr,"  options\n");
	fprintf(stderr,"    -k kind  0/1/2/3/4 which version of BLISS?\n");
	exit(1);
}

const int GEN_KEY = 0;
const int SIGN = 1;
const int VERIFY = 2;

int main(int argc, char *argv[]) {
	int i,c;
	int kind = 0;
	int subcommand = -1;
	entropy_t entropy;
	uint8_t seed[64];
	bliss_param_t params;
	sampler_t sampler;

	/* Parse command line */
	if(argc < 2) {
		usage(argv);
	} else {
		if(strcmp(argv[1],"keygen") == 0) {
			subcommand = GEN_KEY;
			fprintf(stderr,"Generating private key...\n");
		} else if(strcmp(argv[1],"sign") == 0) {
			subcommand = SIGN;
			fprintf(stderr,"Sign...\n");
		} else if(strcmp(argv[1],"verify") == 0) {
			subcommand = VERIFY;
			fprintf(stderr,"Verify...\n");
		} else {
			subcommand = -1;
		}
	}

	while((c = getopt(argc-1,argv+1,"k:")) != -1) {
		switch(c) {
			case 'k':
				kind = atoi(optarg);
				if(kind < 0 || kind > 4)
					usage(argv);
				fprintf(stderr,"Version set to %d\n",kind);
				break;
		}
	}

	/* Initialize entropy */
	for(i = 0; i < 64; i++)
		seed[i] = i%8;
	entropy_init(&entropy,seed);

	/* Initialize parameter */
	bliss_params_init(&params, kind);

	/* Initialize sampler */
	sampler_init(&sampler,params.sigma,params.ell,params.precision,&entropy);

	if(subcommand == GEN_KEY) {
		gen_private_key(&entropy,&params);
	} else if(subcommand == SIGN) {
		test_sign(&entropy,&params);
	} else if(subcommand == VERIFY) {
		test_verify(&entropy,&params);
	} else {
		usage(argv);
	}

	return 0;
}
